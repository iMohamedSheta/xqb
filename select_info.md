```sql
-- Comprehensive MySQL SELECT query showcasing various features
SELECT
    -- Basic column selection with aliases
    e.employee_id AS "Employee ID",
    CONCAT_WS(' ', e.first_name, e.last_name) AS "Full Name",

    -- Arithmetic operations
    (e.salary * 12) AS "Annual Salary",
    ROUND(e.salary * 1.1, 2) AS "Salary After 10% Raise",

    -- Date and time functions
    DATE_FORMAT(e.hire_date, '%M %d, %Y') AS "Formatted Hire Date",
    TIMESTAMPDIFF(YEAR, e.hire_date, CURDATE()) AS "Years of Service",

    -- Conditional logic with CASE
    CASE
        WHEN e.salary < 50000 THEN 'Entry Level'
        WHEN e.salary BETWEEN 50000 AND 100000 THEN 'Mid Level'
        ELSE 'Senior Level'
    END AS "Salary Tier",

    -- Aggregated subquery
    (SELECT AVG(salary) FROM employees WHERE department_id = e.department_id) AS "Dept Avg Salary",

    -- Window functions
    ROW_NUMBER() OVER (PARTITION BY e.department_id ORDER BY e.salary DESC) AS "Salary Rank in Dept",
    DENSE_RANK() OVER (ORDER BY e.salary DESC) AS "Overall Salary Rank",
    LEAD(e.salary, 1, 0) OVER (PARTITION BY e.department_id ORDER BY e.salary) AS "Next Higher Salary",

    -- String functions
    UPPER(d.department_name) AS "Department",
    CHAR_LENGTH(d.department_name) AS "Dept Name Length",

    -- JSON operations (if using MySQL 5.7+)
    JSON_EXTRACT(e.profile_data, '$.skills') AS "Skills",
    JSON_EXTRACT(e.profile_data, '$.education.degree') AS "Degree",

    -- User-defined variable
    @dept_count := (SELECT COUNT(*) FROM departments) AS "Total Departments",

    -- Conditional aggregation
    COUNT(CASE WHEN p.status = 'Completed' THEN 1 END) OVER (PARTITION BY e.employee_id) AS "Completed Projects",

    -- Mathematical functions
    ROUND(SQRT(e.salary), 2) AS "Square Root of Salary",
    POW(2, FLOOR(LOG2(e.salary))) AS "Nearest Power of 2 Below Salary",

    -- Bitwise operations
    e.employee_id & 7 AS "Bitwise AND with 7",

    -- Full-text search relevance score (if FT index exists)
    MATCH(e.bio) AGAINST ('leadership' IN BOOLEAN MODE) AS "Leadership Relevance"

FROM
    employees e

    -- Different types of JOINs
    INNER JOIN departments d ON e.department_id = d.department_id
    LEFT JOIN locations l ON d.location_id = l.location_id
    RIGHT JOIN countries c ON l.country_id = c.country_id
    JOIN regions r ON c.region_id = r.region_id

    -- Self-join
    LEFT JOIN employees m ON e.manager_id = m.employee_id

    -- JOIN with subquery
    JOIN (
        SELECT
            project_id,
            employee_id,
            status,
            RANK() OVER (PARTITION BY employee_id ORDER BY end_date DESC) AS recent_rank
        FROM projects
    ) p ON e.employee_id = p.employee_id AND p.recent_rank = 1

    -- CROSS JOIN for cartesian product (use with caution!)
    CROSS JOIN (SELECT 'Global' as scope) AS scope_table

WHERE
    -- Multiple conditions with different operators
    (e.salary > 50000 OR e.hire_date < '2018-01-01')
    AND d.department_name NOT LIKE 'Admin%'
    AND e.employee_id IN (SELECT employee_id FROM performance WHERE rating >= 4)
    AND EXISTS (SELECT 1 FROM trainings t WHERE t.employee_id = e.employee_id AND t.completed = TRUE)
    AND e.bio IS NOT NULL
    AND (e.email REGEXP '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$')
    AND e.salary BETWEEN 40000 AND 150000

GROUP BY
    e.employee_id,
    d.department_id

HAVING
    COUNT(p.project_id) > 2
    AND MAX(e.salary) < 200000

ORDER BY
    r.region_name ASC,
    d.department_name ASC,
    e.salary DESC

-- Limit and pagination
LIMIT 10 OFFSET 20

-- WITH clause (Common Table Expression - CTE)
WITH department_stats AS (
    SELECT
        department_id,
        COUNT(*) as employee_count,
        AVG(salary) as avg_salary
    FROM employees
    GROUP BY department_id
),
high_performers AS (
    SELECT
        employee_id
    FROM performance
    WHERE rating > 4.5
)

-- UNION, INTERSECT, EXCEPT operations
UNION

-- Second query part of the UNION
SELECT
    e.employee_id AS "Employee ID",
    CONCAT_WS(' ', e.first_name, e.last_name) AS "Full Name",
    -- Include all other columns to match first query's structure
    -- ...
FROM
    former_employees e
    -- Include all necessary joins
    -- ...
WHERE
    e.end_date > DATE_SUB(CURDATE(), INTERVAL 1 YEAR)

-- Procedure/function call within query
CALL update_statistics();

-- Query hints
/*+ INDEX(employees emp_idx) */

```

SELECT Clause Features

Basic column selection with aliases
String manipulation (CONCAT_WS)
Arithmetic operations
Date/time functions
Conditional logic (CASE statements)
Subqueries in the SELECT list
Window functions (ROW_NUMBER, DENSE_RANK, LEAD)
String functions (UPPER, CHAR_LENGTH)
JSON operations
User-defined variables
Conditional aggregation
Mathematical functions
Bitwise operations
Full-text search relevance scoring

FROM Clause Features

Multiple JOIN types (INNER, LEFT, RIGHT, CROSS)
Self joins
Joining with subqueries
Table aliases

WHERE Clause Features

Multiple condition types with different operators
Subqueries (IN, EXISTS)
Pattern matching (LIKE)
Regular expressions
NULL checks
BETWEEN operator

Other Clauses and Features

GROUP BY with multiple columns
HAVING clause with aggregations
ORDER BY with multiple sort criteria
LIMIT and OFFSET for pagination
Common Table Expressions (WITH clause)
Set operations (UNION)
Stored procedure calls
Query hints
