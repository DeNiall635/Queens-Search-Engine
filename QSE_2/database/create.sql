CREATE TABLE IF NOT EXISTS Pages (PageID INT NOT NULL AUTO_INCREMENT, Content text NOT NULL, Url VARCHAR(50), PRIMARY KEY (PageID), INDEX (PageID));
																																		   	
CREATE TABLE IF NOT EXISTS Adverts (AdID INT NOT NULL AUTO_INCREMENT, Keyword varchar(50), Advert text NOT NULL, PRIMARY KEY (AdID));
