CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,     
    room_id VARCHAR(255) NOT NULL,  
    start_time TIMESTAMP NOT NULL,  
    end_time TIMESTAMP NOT NULL     
);

CREATE INDEX idx_room_id ON reservations(room_id);
CREATE INDEX idx_start_time ON reservations(start_time);
