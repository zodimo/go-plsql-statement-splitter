package samples

// GetSimpleSQLSamples returns simple SQL samples for testing
func GetSimpleSQLSamples() map[string]string {
	return map[string]string{
		"simple_select": `SELECT * FROM employees;`,

		"multiple_statements": `
			SELECT * FROM employees;
			DELETE FROM employees WHERE id = 1;
			INSERT INTO employees (id, name) VALUES (1, 'John');
		`,

		"create_procedure": `
			CREATE OR REPLACE PROCEDURE hello_world IS
			BEGIN
				DBMS_OUTPUT.PUT_LINE('Hello, World!');
			END;
			/
		`,

		"anonymous_block": `
			BEGIN
				INSERT INTO employees (id, name) VALUES (1, 'John');
				COMMIT;
			END;
			/
		`,

		"with_comments": `
			-- This is a select statement
			SELECT * FROM employees; -- End of statement
			
			/* This is a multi-line comment
			   with multiple lines */
			DELETE FROM employees WHERE id = 1;
		`,
	}
}

// GetComplexSQLSamples returns complex SQL samples for testing
func GetComplexSQLSamples() map[string]string {
	return map[string]string{
		"package_spec": `
			CREATE OR REPLACE PACKAGE employee_pkg IS
				-- Package constants
				c_default_dept CONSTANT NUMBER := 10;
				
				-- Function to get employee name
				FUNCTION get_employee_name(p_emp_id IN NUMBER) RETURN VARCHAR2;
				
				-- Procedure to update employee
				PROCEDURE update_employee(
					p_emp_id IN NUMBER,
					p_name IN VARCHAR2,
					p_dept_id IN NUMBER DEFAULT c_default_dept
				);
			END employee_pkg;
			/
		`,

		"package_body": `
			CREATE OR REPLACE PACKAGE BODY employee_pkg IS
				-- Private variable
				v_last_updated DATE;
				
				-- Function implementation
				FUNCTION get_employee_name(p_emp_id IN NUMBER) RETURN VARCHAR2 IS
					v_name VARCHAR2(100);
				BEGIN
					SELECT name INTO v_name
					FROM employees
					WHERE id = p_emp_id;
					
					RETURN v_name;
				EXCEPTION
					WHEN NO_DATA_FOUND THEN
						RETURN NULL;
				END get_employee_name;
				
				-- Procedure implementation
				PROCEDURE update_employee(
					p_emp_id IN NUMBER,
					p_name IN VARCHAR2,
					p_dept_id IN NUMBER DEFAULT c_default_dept
				) IS
				BEGIN
					UPDATE employees
					SET name = p_name,
						department_id = p_dept_id
					WHERE id = p_emp_id;
					
					v_last_updated := SYSDATE;
					COMMIT;
				EXCEPTION
					WHEN OTHERS THEN
						ROLLBACK;
						RAISE;
				END update_employee;
				
				-- Initialize
				BEGIN
					v_last_updated := SYSDATE;
				END employee_pkg;
			/
		`,

		"nested_blocks": `
			BEGIN
				FOR r_emp IN (SELECT * FROM employees) LOOP
					BEGIN
						UPDATE departments
						SET employee_count = employee_count + 1
						WHERE id = r_emp.department_id;
					EXCEPTION
						WHEN NO_DATA_FOUND THEN
							INSERT INTO departments (id, name, employee_count)
							VALUES (r_emp.department_id, 'New Department', 1);
					END;
				END LOOP;
				COMMIT;
			END;
			/
		`,
	}
}

// GetInvalidSQLSamples returns invalid SQL samples for testing error handling
func GetInvalidSQLSamples() map[string]string {
	return map[string]string{
		"missing_semicolon": `
			SELECT * FROM employees
			DELETE FROM employees WHERE id = 1;
		`,

		"unclosed_comment": `
			SELECT * FROM employees;
			/* This comment is not closed
			DELETE FROM employees WHERE id = 1;
		`,

		"unclosed_block": `
			BEGIN
				DBMS_OUTPUT.PUT_LINE('Hello, World!');
			-- Missing END statement
			/
		`,
	}
}
