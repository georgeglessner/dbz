-- Seed data for blog database
-- Sample users, posts, and comments

-- Insert sample users
INSERT INTO users (username, email, password_hash, full_name, bio) VALUES
('johndoe', 'john@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'John Doe', 'Software engineer and tech blogger'),
('janedoe', 'jane@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Jane Smith', 'Full-stack developer and open source enthusiast'),
('bobwilson', 'bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Bob Wilson', 'DevOps engineer and cloud architect')
ON CONFLICT (email) DO NOTHING;

-- Insert categories
INSERT INTO categories (name, slug, description) VALUES
('Technology', 'technology', 'Latest tech news and tutorials'),
('Programming', 'programming', 'Coding tips and best practices'),
('DevOps', 'devops', 'Infrastructure and deployment guides'),
('Database', 'database', 'Database design and optimization')
ON CONFLICT (slug) DO NOTHING;

-- Insert tags
INSERT INTO tags (name, slug) VALUES
('Python', 'python'),
('Go', 'go'),
('PostgreSQL', 'postgresql'),
('Docker', 'docker'),
('Kubernetes', 'kubernetes'),
('Tutorial', 'tutorial'),
('Best Practices', 'best-practices')
ON CONFLICT (slug) DO NOTHING;

-- Insert sample posts
INSERT INTO posts (user_id, title, slug, content, status, published_at, category_id) VALUES
(1, 'Getting Started with PostgreSQL', 'getting-started-with-postgresql', 
'PostgreSQL is a powerful, open source object-relational database system. In this post, we will cover the basics of setting up PostgreSQL and performing common operations.

Installation

PostgreSQL can be installed on various platforms. For macOS users, the easiest way is using Homebrew. For Linux users, use your package manager. Windows users can use the installer from the PostgreSQL website.

Creating Your First Database

Once installed, you can create a database using the createdb command or by connecting to the default postgres database and running CREATE DATABASE.

Basic Operations

Here are some essential SQL commands to get you started: creating tables, inserting data, and querying. These operations form the foundation of working with any relational database.

Stay tuned for more advanced PostgreSQL topics!',
'published', NOW() - INTERVAL '5 days', 4),

(2, 'Introduction to Go Programming', 'introduction-to-go-programming',
'Go, also known as Golang, is a statically typed, compiled programming language designed at Google. It is known for its simplicity, efficiency, and great concurrency support.

Why Go?

Go offers simple syntax that is easy to learn and read. It has fast compilation even for large projects. The language provides great concurrency support through goroutines and channels. It comes with a strong standard library and supports cross-compilation for multiple platforms.

Hello World

Getting started with Go is straightforward. You write your code, build it with go build, and run the resulting binary. The tooling is excellent and included in the standard distribution.

Concurrency Example

One of Go standout features is its approach to concurrency. Goroutines make it easy to run functions concurrently, and channels provide a safe way to communicate between them.

Go is an excellent choice for building scalable network services and distributed systems.',
'published', NOW() - INTERVAL '3 days', 2),

(3, 'Docker Best Practices for Production', 'docker-best-practices-for-production',
'Docker has revolutionized how we deploy applications. Here are some best practices for running Docker in production environments.

1. Use Official Base Images

Always start with official images from Docker Hub. They are maintained, regularly updated, and follow security best practices. Avoid using random user images for production workloads.

2. Minimize Image Size

Smaller images mean faster deployments and less attack surface. Use multi-stage builds to compile your application in one stage and copy only the artifacts to a smaller runtime image.

3. Do Not Run as Root

Create a non-root user for your application. This is a security best practice that limits the damage if your application is compromised.

4. Use dockerignore

Prevent unnecessary files from being added to your image. Exclude node_modules, git directories, environment files, and log files.

5. Health Checks

Add health checks to your containers so Docker can monitor the application status and restart unhealthy containers automatically.

Following these practices will help you build more secure, efficient, and maintainable Docker deployments.',
'published', NOW() - INTERVAL '1 day', 3)
ON CONFLICT (slug) DO NOTHING;

-- Add tags to posts
INSERT INTO post_tags (post_id, tag_id) VALUES
(1, 3), (1, 6), -- PostgreSQL post: postgresql, tutorial
(2, 2), (2, 7), -- Go post: go, best-practices
(3, 4), (3, 7)  -- Docker post: docker, best-practices
ON CONFLICT DO NOTHING;

-- Insert sample comments
INSERT INTO comments (post_id, user_id, content) VALUES
(1, 2, 'Great introduction! Really helped me get started with PostgreSQL.'),
(1, 3, 'Could you add a section about indexing? That would be really helpful.'),
(2, 1, 'Go is such a great language. The concurrency model is really elegant.'),
(2, 3, 'I love how fast Go compiles. Makes development so much smoother.'),
(3, 1, 'These are solid practices. We implemented multi-stage builds and saw 40 percent reduction in image size.'),
(3, 2, 'The multi-stage build example is exactly what I needed. Thanks!')
ON CONFLICT DO NOTHING;
