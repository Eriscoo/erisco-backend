DELETE FROM users WHERE email = 'eriscobert@gmail.com';
DELETE FROM tags WHERE name IN ('ubuntu', 'golang', 'react', 'javascript', 'tutorial', 'technology', 'website', 'typescript', 'python', 'java', 'rust', 'vue', 'nodejs', 'docker', 'kubernetes', 'postgresql', 'mysql', 'redis', 'graphql', 'rest-api', 'testing', 'devops', 'security', 'performance', 'database', 'frontend', 'backend', 'mobile', 'erisco');
DELETE FROM categories WHERE name IN ('general', 'technology', 'programming', 'design', 'business', 'lifestyle', 'health', 'science', 'education', 'travel', 'food', 'music', 'art', 'photography', 'sports', 'gaming', 'news', 'finance', 'marketing', 'startup', 'career', 'tutorial', 'review', 'opinion', 'interview', 'case-study');
DELETE FROM user_role WHERE role_name IN ('super admin');
