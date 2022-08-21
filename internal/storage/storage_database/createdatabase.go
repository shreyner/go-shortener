package storagedatabase

import (
	"context"
	"time"
)

func (s *storageSQL) CheckAndCreateSchema() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, `
		create table if not exists short_url
		(
			id         		varchar                   not null,
			url        		varchar                   not null,
			user_id    		varchar,
			created_at 		date default current_date not null,
			correlation_id 	varchar,
			deleted 		boolean default false     not null
		);
		
		create unique index if not exists short_url_id_uindex
			on short_url (id);
		
		create unique index if not exists short_url_uindex
		    on short_url (url);
	`)

	return err
}
