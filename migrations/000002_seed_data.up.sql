INSERT INTO user_role (role_name) VALUES
    ('super admin')
ON CONFLICT (role_name) DO NOTHING;

INSERT INTO users (name, email, password_hash, role_id) VALUES
    ('Erisco Berto', 'eriscobert@gmail.com', '$2a$10$8SuSEwuvSIIgahDBq4WBJe2Szo9BUUMfYQ06vuuZBRnsp22c/eroC', 1)
ON CONFLICT (email) DO NOTHING;

INSERT INTO categories (name) VALUES
    ('general'),
    ('technology'),
    ('programming'),
    ('design'),
    ('business'),
    ('lifestyle'),
    ('health'),
    ('science'),
    ('education'),
    ('travel'),
    ('food'),
    ('music'),
    ('art'),
    ('photography'),
    ('sports'),
    ('gaming'),
    ('news'),
    ('finance'),
    ('marketing'),
    ('startup'),
    ('career'),
    ('tutorial'),
    ('review'),
    ('opinion'),
    ('interview'),
    ('case-study')
ON CONFLICT (name) DO NOTHING;

INSERT INTO tags (name) VALUES
    ('ubuntu'),
    ('golang'),
    ('react'),
    ('javascript'),
    ('tutorial'),
    ('technology'),
    ('website'),
    ('typescript'),
    ('python'),
    ('java'),
    ('rust'),
    ('vue'),
    ('nodejs'),
    ('docker'),
    ('kubernetes'),
    ('postgresql'),
    ('mysql'),
    ('redis'),
    ('graphql'),
    ('rest-api'),
    ('testing'),
    ('devops'),
    ('security'),
    ('performance'),
    ('database'),
    ('frontend'),
    ('backend'),
    ('mobile'),
    ('erisco')
ON CONFLICT (name) DO NOTHING;
