-- Seed data for CVWO Forum
-- Run this after migration.sql

-- Clear existing data (optional - comment out if you want to preserve existing data)
TRUNCATE TABLE comments, posts, topics, users RESTART IDENTITY CASCADE;

-- Insert Users
INSERT INTO users (username, created_at) VALUES
    ('alice', NOW() - INTERVAL '30 days'),
    ('bob', NOW() - INTERVAL '25 days'),
    ('charlie', NOW() - INTERVAL '20 days'),
    ('diana', NOW() - INTERVAL '15 days'),
    ('eve', NOW() - INTERVAL '10 days');

-- Insert Topics
INSERT INTO topics (title, description, user_id, created_at, updated_at) VALUES
    ('General Discussion', 'A place to discuss anything and everything', 1, NOW() - INTERVAL '28 days', NOW() - INTERVAL '28 days'),
    ('Programming Help', 'Get help with your coding questions', 2, NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
    ('Project Showcase', 'Share your projects with the community', 1, NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),
    ('Off-Topic', 'Random discussions that don''t fit elsewhere', 3, NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
    ('Announcements', 'Important updates and announcements', 1, NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days');

-- Insert Posts
INSERT INTO posts (title, content, topic_id, user_id, created_at, updated_at) VALUES
    -- General Discussion posts
    ('Welcome to the forum!', 'Hey everyone! Welcome to our new forum. Feel free to introduce yourself here.', 1, 1, NOW() - INTERVAL '27 days', NOW() - INTERVAL '27 days'),
    ('What are you working on?', 'I''m curious what projects everyone is currently working on. Share your progress!', 1, 2, NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
    ('Forum rules and guidelines', 'Please be respectful to all members. No spam or self-promotion without permission.', 1, 1, NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),
    
    -- Programming Help posts
    ('How to learn Go?', 'I''m new to Go programming. What resources would you recommend for beginners?', 2, 3, NOW() - INTERVAL '23 days', NOW() - INTERVAL '23 days'),
    ('React vs Vue - which one?', 'I''m starting a new frontend project. Should I use React or Vue? What are the pros and cons?', 2, 4, NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
    ('Database design question', 'I''m designing a database for an e-commerce app. Should I use SQL or NoSQL?', 2, 5, NOW() - INTERVAL '19 days', NOW() - INTERVAL '19 days'),
    ('Help with TypeScript generics', 'Can someone explain TypeScript generics with a simple example? I''m confused about the syntax.', 2, 3, NOW() - INTERVAL '17 days', NOW() - INTERVAL '17 days'),
    
    -- Project Showcase posts
    ('My first Go API', 'Just finished building my first REST API in Go! It''s a simple todo app but I learned a lot.', 3, 2, NOW() - INTERVAL '15 days', NOW() - INTERVAL '15 days'),
    ('Portfolio website redesign', 'Check out my new portfolio website! Built with React and TailwindCSS. Feedback welcome!', 3, 4, NOW() - INTERVAL '13 days', NOW() - INTERVAL '13 days'),
    ('Open source CLI tool', 'I created a CLI tool for managing dotfiles. It''s open source, contributions welcome!', 3, 5, NOW() - INTERVAL '11 days', NOW() - INTERVAL '11 days'),
    
    -- Off-Topic posts
    ('Favorite programming music?', 'What do you listen to while coding? I need some new playlist recommendations!', 4, 3, NOW() - INTERVAL '9 days', NOW() - INTERVAL '9 days'),
    ('Mechanical keyboards', 'Anyone here into mechanical keyboards? Just got my first one and I''m hooked!', 4, 2, NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
    
    -- Announcements posts
    ('Forum maintenance scheduled', 'The forum will be down for maintenance this Saturday from 2-4 AM UTC.', 5, 1, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
    ('New features coming soon', 'We''re working on adding dark mode and notification features. Stay tuned!', 5, 1, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days');

-- Insert Comments
INSERT INTO comments (content, post_id, user_id, created_at, updated_at) VALUES
    -- Comments on "Welcome to the forum!"
    ('Thanks for setting this up! Excited to be here.', 1, 2, NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
    ('Hello everyone! I''m Bob, a backend developer from NYC.', 1, 2, NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
    ('Hi! I''m Charlie, learning web development. Nice to meet you all!', 1, 3, NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
    ('Great to see a new community forming. Looking forward to discussions!', 1, 4, NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
    
    -- Comments on "How to learn Go?"
    ('I recommend "A Tour of Go" on the official website. It''s interactive and free!', 4, 1, NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
    ('The Go documentation is excellent. Also check out Exercism for practice problems.', 4, 2, NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
    ('Thanks for the suggestions! I''ll check those out.', 4, 3, NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
    ('"Learning Go" by Jon Bodner is a great book if you prefer reading.', 4, 5, NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
    
    -- Comments on "React vs Vue"
    ('I prefer React because of the larger ecosystem and job market.', 5, 1, NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
    ('Vue has a gentler learning curve. Great for beginners!', 5, 3, NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
    ('Both are great choices. Pick one and stick with it!', 5, 5, NOW() - INTERVAL '19 days', NOW() - INTERVAL '19 days'),
    
    -- Comments on "My first Go API"
    ('Congrats! Building your first API is a big milestone.', 8, 1, NOW() - INTERVAL '14 days', NOW() - INTERVAL '14 days'),
    ('Nice work! What did you use for the database?', 8, 3, NOW() - INTERVAL '14 days', NOW() - INTERVAL '14 days'),
    ('I used PostgreSQL with the standard library. Worked great!', 8, 2, NOW() - INTERVAL '13 days', NOW() - INTERVAL '13 days'),
    
    -- Comments on "Favorite programming music?"
    ('Lo-fi hip hop all the way! Those YouTube livestreams are perfect.', 11, 2, NOW() - INTERVAL '8 days', NOW() - INTERVAL '8 days'),
    ('I prefer silence or white noise. Music distracts me.', 11, 4, NOW() - INTERVAL '8 days', NOW() - INTERVAL '8 days'),
    ('Video game soundtracks! They''re designed to help you focus.', 11, 5, NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
    ('Try "musicForProgramming()" - it''s made specifically for coding!', 11, 1, NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
    
    -- Comments on "Mechanical keyboards"
    ('Welcome to the rabbit hole! Your wallet will never recover.', 12, 4, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
    ('What switches did you get? I''m a big fan of Cherry MX Browns.', 12, 5, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
    ('I got Gateron Reds. Super smooth for typing!', 12, 2, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
    
    -- Comments on "New features coming soon"
    ('Dark mode! Finally! Can''t wait.', 14, 2, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
    ('Will there be email notifications?', 14, 3, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
    ('Yes, email notifications are planned for the next release!', 14, 1, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day');

-- Verify seed data
SELECT 'Users:' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Topics:', COUNT(*) FROM topics
UNION ALL
SELECT 'Posts:', COUNT(*) FROM posts
UNION ALL
SELECT 'Comments:', COUNT(*) FROM comments;