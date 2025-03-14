-- Add a test admin user
INSERT INTO users (email, is_admin) VALUES
('test.admin@example.com', true);

-- Add two example events
INSERT INTO events (
    title,
    description,
    image_url,
    system,
    host_name,
    email,
    phone_number,
    host,
    pulje_name,
    max_players,
    child_friendly,
    adults_only,
    beginner_friendly,
    experienced_only,
    can_be_run_in_english,
    long_running,
    short_running,
    status
) VALUES (
    'Dungeons & Dragons: Den Tapte Minen',
    'Bli med på et spennende eventyr i den klassiske D&D-modulen "Den Tapte Minen av Phandelver". Perfekt for nye spillere!',
    'https://imgur.com/example1',
    'D&D 5e',
    'Erik Spilleder',
    'test.admin@example.com',
    12345678,
    1, -- refererer til test.admin@example.com
    'Lørdag morgen',
    6,
    false,
    false,
    true,
    false,
    true,
    false,
    true,
    'Publisert'
),
(
    'Vampire: Nattens Barn',
    'En intens fortelling om intriger og makt i Oslos vampyrsamfunn. Kun for erfarne rollespillere.',
    'https://imgur.com/example2',
    'Vampire: The Masquerade 5th Edition',
    'Maria Storyteller',
    'test.admin@example.com',
    12345678,
    1, -- refererer til test.admin@example.com
    'Lørdag kveld',
    4,
    false,
    true,
    false,
    true,
    false,
    true,
    false,
    'Publisert'
);
