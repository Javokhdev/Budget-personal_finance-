CREATE TABLE country(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name varchar(255) not NULL,
  flag varchar(255) not NULL,
  total INTEGER DEFAULT 0
);

CREATE TABLE team (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,                     
    country_id UUID NOT NULL REFERENCES country(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
);

CREATE TABLE athlete (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL ,                     
    country_id UUID NOT NULL REFERENCES country(id),    
    team_id uuid REFERENCES team(id),                  
    sport_type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0          
);

CREATE TABLE event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,                     
    sport_type VARCHAR(100) NOT NULL,              
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    outcome VARCHAR not null,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at BIGINT DEFAULT 0
);

create type typenum as enum ('gold', 'silver', 'bronze');

CREATE TABLE medals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  country_id UUID REFERENCES country(id) not NULL,
  event_id UUID REFERENCES event(id) not NULL,
  athlete_id UUID REFERENCES athlete(id) not NULL,
  type typenum not NULL, 
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at bigint DEFAULT 0
);
create Table eventparticipant (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID not null REFERENCES event(id),
  athlete_id UUID REFERENCES athlete(id),
  team_id UUID REFERENCES team(id),
  created_at TIMESTAMP DEFAULT NOW(),
  deleted_at BIGINT DEFAULT 0
);