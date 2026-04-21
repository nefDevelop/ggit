# Table "access"

```
 Field  | Column  |   PostgreSQL    |        SQLite3        
--------+---------+-----------------+-----------------------
 ID     | id      | BIGSERIAL       | INTEGER AUTOINCREMENT 
 UserID | user_id | BIGINT NOT NULL | INTEGER NOT NULL      
 RepoID | repo_id | BIGINT NOT NULL | INTEGER NOT NULL      
 Mode   | mode    | BIGINT NOT NULL | INTEGER NOT NULL      

Primary keys: id
Indexes: 
	"access_user_repo_unique" UNIQUE (user_id, repo_id)
```

# Table "access_token"

```
    Field    |    Column    |         PostgreSQL          |           SQLite3           
-------------+--------------+-----------------------------+-----------------------------
 ID          | id           | BIGSERIAL                   | INTEGER AUTOINCREMENT       
 UserID      | uid          | BIGINT                      | INTEGER                     
 Name        | name         | TEXT                        | TEXT                        
 Sha1        | sha1         | VARCHAR(40) UNIQUE          | VARCHAR(40) UNIQUE          
 SHA256      | sha256       | VARCHAR(64) NOT NULL UNIQUE | VARCHAR(64) NOT NULL UNIQUE 
 CreatedUnix | created_unix | BIGINT                      | INTEGER                     
 UpdatedUnix | updated_unix | BIGINT                      | INTEGER                     

Primary keys: id
Indexes: 
	"idx_access_token_user_id" (uid)
```

# Table "action"

```
    Field     |     Column     |           PostgreSQL           |            SQLite3             
--------------+----------------+--------------------------------+--------------------------------
 ID           | id             | BIGSERIAL                      | INTEGER AUTOINCREMENT          
 UserID       | user_id        | BIGINT                         | INTEGER                        
 OpType       | op_type        | BIGINT                         | INTEGER                        
 ActUserID    | act_user_id    | BIGINT                         | INTEGER                        
 ActUserName  | act_user_name  | TEXT                           | TEXT                           
 RepoID       | repo_id        | BIGINT                         | INTEGER                        
 RepoUserName | repo_user_name | TEXT                           | TEXT                           
 RepoName     | repo_name      | TEXT                           | TEXT                           
 RefName      | ref_name       | TEXT                           | TEXT                           
 IsPrivate    | is_private     | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 Content      | content        | TEXT                           | TEXT                           
 CreatedUnix  | created_unix   | BIGINT                         | INTEGER                        

Primary keys: id
Indexes: 
	"idx_action_repo_id" (repo_id)
	"idx_action_user_id" (user_id)
```

# Table "attachment"

```
    Field    |    Column    |      PostgreSQL      |        SQLite3        
-------------+--------------+----------------------+-----------------------
 ID          | id           | BIGSERIAL            | INTEGER AUTOINCREMENT 
 UUID        | uuid         | TEXT NOT NULL UNIQUE | TEXT NOT NULL UNIQUE  
 IssueID     | issue_id     | BIGINT               | INTEGER               
 CommentID   | comment_id   | BIGINT               | INTEGER               
 ReleaseID   | release_id   | BIGINT               | INTEGER               
 Name        | name         | TEXT NOT NULL        | TEXT NOT NULL         
 CreatedUnix | created_unix | BIGINT               | INTEGER               

Primary keys: id
Indexes: 
	"idx_attachment_comment_id" (comment_id)
	"idx_attachment_issue_id" (issue_id)
	"idx_attachment_release_id" (release_id)
```

# Table "collaboration"

```
 Field  | Column  |        PostgreSQL         |          SQLite3           
--------+---------+---------------------------+----------------------------
 ID     | id      | BIGSERIAL                 | INTEGER AUTOINCREMENT      
 UserID | user_id | BIGINT NOT NULL           | INTEGER NOT NULL           
 RepoID | repo_id | BIGINT NOT NULL           | INTEGER NOT NULL           
 Mode   | mode    | BIGINT NOT NULL DEFAULT 2 | INTEGER NOT NULL DEFAULT 2 

Primary keys: id
Indexes: 
	"collaboration_user_repo_unique" UNIQUE (user_id, repo_id)
	"idx_collaboration_repo_id" (repo_id)
	"idx_collaboration_user_id" (user_id)
```

# Table "comment"

```
    Field    |    Column    |   PostgreSQL    |        SQLite3        
-------------+--------------+-----------------+-----------------------
 ID          | id           | BIGSERIAL       | INTEGER AUTOINCREMENT 
 Type        | type         | BIGINT NOT NULL | INTEGER NOT NULL      
 PosterID    | poster_id    | BIGINT          | INTEGER               
 IssueID     | issue_id     | BIGINT          | INTEGER               
 CommitID    | commit_id    | BIGINT          | INTEGER               
 Line        | line         | BIGINT          | INTEGER               
 Content     | content      | TEXT            | TEXT                  
 CreatedUnix | created_unix | BIGINT          | INTEGER               
 UpdatedUnix | updated_unix | BIGINT          | INTEGER               
 CommitSHA   | commit_sha   | VARCHAR(40)     | VARCHAR(40)           

Primary keys: id
Indexes: 
	"idx_comment_issue_id" (issue_id)
	"idx_comment_poster_id" (poster_id)
```

# Table "deploy_key"

```
    Field    |    Column    |        PostgreSQL         |          SQLite3           
-------------+--------------+---------------------------+----------------------------
 ID          | id           | BIGSERIAL                 | INTEGER AUTOINCREMENT      
 KeyID       | key_id       | BIGINT                    | INTEGER                    
 RepoID      | repo_id      | BIGINT                    | INTEGER                    
 Name        | name         | TEXT NOT NULL             | TEXT NOT NULL              
 Fingerprint | fingerprint  | TEXT NOT NULL             | TEXT NOT NULL              
 Mode        | mode         | BIGINT NOT NULL DEFAULT 2 | INTEGER NOT NULL DEFAULT 2 
 CreatedUnix | created_unix | BIGINT                    | INTEGER                    
 UpdatedUnix | updated_unix | BIGINT                    | INTEGER                    

Primary keys: id
Indexes: 
	"deploy_key_key_repo_unique" UNIQUE (key_id, repo_id)
	"idx_deploy_key_key_id" (key_id)
	"idx_deploy_key_repo_id" (repo_id)
```

# Table "email_address"

```
    Field    |    Column    |           PostgreSQL           |            SQLite3             
-------------+--------------+--------------------------------+--------------------------------
 ID          | id           | BIGSERIAL                      | INTEGER AUTOINCREMENT          
 UserID      | uid          | BIGINT NOT NULL                | INTEGER NOT NULL               
 Email       | email        | VARCHAR(254) NOT NULL          | TEXT NOT NULL                  
 IsActivated | is_activated | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 

Primary keys: id
Indexes: 
	"email_address_user_email_unique" UNIQUE (uid, email)
	"idx_email_address_user_id" (uid)
```

# Table "follow"

```
  Field   |  Column   |   PostgreSQL    |        SQLite3        
----------+-----------+-----------------+-----------------------
 ID       | id        | BIGSERIAL       | INTEGER AUTOINCREMENT 
 UserID   | user_id   | BIGINT NOT NULL | INTEGER NOT NULL      
 FollowID | follow_id | BIGINT NOT NULL | INTEGER NOT NULL      

Primary keys: id
Indexes: 
	"follow_user_follow_unique" UNIQUE (user_id, follow_id)
```

# Table "hook_task"

```
      Field      |      Column      |         PostgreSQL          |        SQLite3        
-----------------+------------------+-----------------------------+-----------------------
 ID              | id               | BIGSERIAL                   | INTEGER AUTOINCREMENT 
 RepoID          | repo_id          | BIGINT                      | INTEGER               
 HookID          | hook_id          | BIGINT                      | INTEGER               
 UUID            | uuid             | VARCHAR(36) NOT NULL UNIQUE | TEXT NOT NULL UNIQUE  
 Type            | type             | BIGINT                      | INTEGER               
 URL             | url              | TEXT NOT NULL               | TEXT NOT NULL         
 Signature       | signature        | TEXT                        | TEXT                  
 PayloadContent  | payload_content  | TEXT                        | TEXT                  
 ContentType     | content_type     | BIGINT                      | INTEGER               
 EventType       | event_type       | TEXT                        | TEXT                  
 IsSSL           | is_ssl           | BOOLEAN                     | NUMERIC               
 IsDelivered     | is_delivered     | BOOLEAN                     | NUMERIC               
 Delivered       | delivered        | BIGINT                      | INTEGER               
 IsSucceed       | is_succeed       | BOOLEAN                     | NUMERIC               
 RequestContent  | request_content  | TEXT                        | TEXT                  
 ResponseContent | response_content | TEXT                        | TEXT                  

Primary keys: id
Indexes: 
	"idx_hook_task_hook_id" (hook_id)
	"idx_hook_task_repo_id" (repo_id)
```

# Table "issue"

```
    Field     |    Column     |   PostgreSQL    |        SQLite3        
--------------+---------------+-----------------+-----------------------
 ID           | id            | BIGSERIAL       | INTEGER AUTOINCREMENT 
 RepoID       | repo_id       | BIGINT NOT NULL | INTEGER NOT NULL      
 Index        | index         | BIGINT NOT NULL | INTEGER NOT NULL      
 PosterID     | poster_id     | BIGINT          | INTEGER               
 Title        | name          | TEXT NOT NULL   | TEXT NOT NULL         
 Content      | content       | TEXT            | TEXT                  
 MilestoneID  | milestone_id  | BIGINT          | INTEGER               
 Priority     | priority      | BIGINT          | INTEGER               
 AssigneeID   | assignee_id   | BIGINT          | INTEGER               
 IsClosed     | is_closed     | BOOLEAN         | NUMERIC               
 IsPull       | is_pull       | BOOLEAN         | NUMERIC               
 NumComments  | num_comments  | BIGINT          | INTEGER               
 DeadlineUnix | deadline_unix | BIGINT          | INTEGER               
 CreatedUnix  | created_unix  | BIGINT          | INTEGER               
 UpdatedUnix  | updated_unix  | BIGINT          | INTEGER               

Primary keys: id
Indexes: 
	"idx_issue_assignee_id" (assignee_id)
	"idx_issue_milestone_id" (milestone_id)
	"idx_issue_poster_id" (poster_id)
	"idx_issue_repo_id" (repo_id)
	"issue_repo_index_unique" UNIQUE (repo_id, index)
```

# Table "issue_label"

```
  Field  |  Column  | PostgreSQL |        SQLite3        
---------+----------+------------+-----------------------
 ID      | id       | BIGSERIAL  | INTEGER AUTOINCREMENT 
 IssueID | issue_id | BIGINT     | INTEGER               
 LabelID | label_id | BIGINT     | INTEGER               

Primary keys: id
Indexes: 
	"issue_label_issue_label_unique" UNIQUE (issue_id, label_id)
```

# Table "issue_user"

```
    Field    |    Column    | PostgreSQL |        SQLite3        
-------------+--------------+------------+-----------------------
 ID          | id           | BIGSERIAL  | INTEGER AUTOINCREMENT 
 UserID      | uid          | BIGINT     | INTEGER               
 IssueID     | issue_id     | BIGINT     | INTEGER               
 RepoID      | repo_id      | BIGINT     | INTEGER               
 MilestoneID | milestone_id | BIGINT     | INTEGER               
 IsRead      | is_read      | BOOLEAN    | NUMERIC               
 IsAssigned  | is_assigned  | BOOLEAN    | NUMERIC               
 IsMentioned | is_mentioned | BOOLEAN    | NUMERIC               
 IsPoster    | is_poster    | BOOLEAN    | NUMERIC               
 IsClosed    | is_closed    | BOOLEAN    | NUMERIC               

Primary keys: id
Indexes: 
	"idx_issue_user_repo_id" (repo_id)
	"idx_issue_user_user_id" (uid)
```

# Table "label"

```
      Field      |      Column       |     PostgreSQL      |        SQLite3        
-----------------+-------------------+---------------------+-----------------------
 ID              | id                | BIGSERIAL           | INTEGER AUTOINCREMENT 
 RepoID          | repo_id           | BIGINT              | INTEGER               
 Name            | name              | TEXT NOT NULL       | TEXT NOT NULL         
 Color           | color             | VARCHAR(7) NOT NULL | VARCHAR(7) NOT NULL   
 NumIssues       | num_issues        | BIGINT              | INTEGER               
 NumClosedIssues | num_closed_issues | BIGINT              | INTEGER               

Primary keys: id
Indexes: 
	"idx_label_repo_id" (repo_id)
```

# Table "lfs_object"

```
    Field    |    Column    |   PostgreSQL    |        SQLite3        
-------------+--------------+-----------------+-----------------------
 ID          | id           | BIGSERIAL       | INTEGER AUTOINCREMENT 
 OID         | oid          | TEXT NOT NULL   | TEXT NOT NULL         
 Size        | size         | BIGINT NOT NULL | INTEGER NOT NULL      
 RepoID      | repo_id      | BIGINT NOT NULL | INTEGER NOT NULL      
 Storage     | storage      | TEXT NOT NULL   | TEXT NOT NULL         
 CreatedUnix | created_unix | BIGINT          | INTEGER               

Primary keys: id
Indexes: 
	"idx_lfs_object_o_id" (oid)
	"lfs_object_oid_repo_unique" UNIQUE (oid, repo_id)
```

# Table "login_source"

```
    Field    |    Column    |    PostgreSQL    |        SQLite3        
-------------+--------------+------------------+-----------------------
 ID          | id           | BIGSERIAL        | INTEGER AUTOINCREMENT 
 Type        | type         | BIGINT           | INTEGER               
 Name        | name         | TEXT UNIQUE      | TEXT UNIQUE           
 IsActived   | is_actived   | BOOLEAN NOT NULL | NUMERIC NOT NULL      
 IsDefault   | is_default   | BOOLEAN          | NUMERIC               
 Config      | cfg          | TEXT             | TEXT                  
 CreatedUnix | created_unix | BIGINT           | INTEGER               
 UpdatedUnix | updated_unix | BIGINT           | INTEGER               

Primary keys: id
```

# Table "milestone"

```
      Field      |      Column       |  PostgreSQL   |        SQLite3        
-----------------+-------------------+---------------+-----------------------
 ID              | id                | BIGSERIAL     | INTEGER AUTOINCREMENT 
 RepoID          | repo_id           | BIGINT        | INTEGER               
 Name            | name              | TEXT NOT NULL | TEXT NOT NULL         
 Content         | content           | TEXT          | TEXT                  
 IsClosed        | is_closed         | BOOLEAN       | NUMERIC               
 NumIssues       | num_issues        | BIGINT        | INTEGER               
 NumClosedIssues | num_closed_issues | BIGINT        | INTEGER               
 Completeness    | completeness      | BIGINT        | INTEGER               
 DeadlineUnix    | deadline_unix     | BIGINT        | INTEGER               
 ClosedDateUnix  | closed_date_unix  | BIGINT        | INTEGER               

Primary keys: id
Indexes: 
	"idx_milestone_repo_id" (repo_id)
```

# Table "mirror"

```
    Field     |      Column      |          PostgreSQL           |            SQLite3            
--------------+------------------+-------------------------------+-------------------------------
 ID           | id               | BIGSERIAL                     | INTEGER AUTOINCREMENT         
 RepoID       | repo_id          | BIGINT                        | INTEGER                       
 Interval     | interval         | BIGINT                        | INTEGER                       
 EnablePrune  | enable_prune     | BOOLEAN NOT NULL DEFAULT TRUE | NUMERIC NOT NULL DEFAULT TRUE 
 LastSyncUnix | updated_unix     | BIGINT                        | INTEGER                       
 NextSyncUnix | next_update_unix | BIGINT                        | INTEGER                       

Primary keys: id
Indexes: 
	"idx_mirror_repo_id" (repo_id)
```

# Table "notice"

```
    Field    |    Column    |   PostgreSQL    |        SQLite3        
-------------+--------------+-----------------+-----------------------
 ID          | id           | BIGSERIAL       | INTEGER AUTOINCREMENT 
 Type        | type         | BIGINT NOT NULL | INTEGER NOT NULL      
 Description | description  | TEXT            | TEXT                  
 CreatedUnix | created_unix | BIGINT          | INTEGER               

Primary keys: id
```

# Table "org_user"

```
  Field   |  Column   |           PostgreSQL           |            SQLite3             
----------+-----------+--------------------------------+--------------------------------
 ID       | id        | BIGSERIAL                      | INTEGER AUTOINCREMENT          
 UID      | uid       | BIGINT NOT NULL                | INTEGER NOT NULL               
 OrgID    | org_id    | BIGINT NOT NULL                | INTEGER NOT NULL               
 IsPublic | is_public | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 IsOwner  | is_owner  | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 NumTeams | num_teams | BIGINT NOT NULL DEFAULT 0      | INTEGER NOT NULL DEFAULT 0     

Primary keys: id
Indexes: 
	"idx_org_user_org_id" (org_id)
	"idx_org_user_uid" (uid)
	"org_user_user_org_unique" UNIQUE (uid, org_id)
```

# Table "protect_branch"

```
       Field        |        Column        |   PostgreSQL    |        SQLite3        
--------------------+----------------------+-----------------+-----------------------
 ID                 | id                   | BIGSERIAL       | INTEGER AUTOINCREMENT 
 RepoID             | repo_id              | BIGINT NOT NULL | INTEGER NOT NULL      
 Name               | name                 | TEXT NOT NULL   | TEXT NOT NULL         
 Protected          | protected            | BOOLEAN         | NUMERIC               
 RequirePullRequest | require_pull_request | BOOLEAN         | NUMERIC               
 EnableWhitelist    | enable_whitelist     | BOOLEAN         | NUMERIC               
 WhitelistUserIDs   | whitelist_user_ids   | TEXT            | TEXT                  
 WhitelistTeamIDs   | whitelist_team_ids   | TEXT            | TEXT                  

Primary keys: id
Indexes: 
	"protect_branch_unique" UNIQUE (repo_id, name)
```

# Table "protect_branch_whitelist"

```
      Field      |      Column       | PostgreSQL |        SQLite3        
-----------------+-------------------+------------+-----------------------
 ID              | id                | BIGSERIAL  | INTEGER AUTOINCREMENT 
 ProtectBranchID | protect_branch_id | BIGINT     | INTEGER               
 RepoID          | repo_id           | BIGINT     | INTEGER               
 Name            | name              | TEXT       | TEXT                  
 UserID          | user_id           | BIGINT     | INTEGER               

Primary keys: id
Indexes: 
	"idx_protect_branch_whitelist_protect_branch_id" (protect_branch_id)
	"protect_branch_whitelist_unique" UNIQUE (repo_id, name, user_id)
```

# Table "public_key"

```
    Field    |    Column    |        PostgreSQL         |          SQLite3           
-------------+--------------+---------------------------+----------------------------
 ID          | id           | BIGSERIAL                 | INTEGER AUTOINCREMENT      
 OwnerID     | owner_id     | BIGINT NOT NULL           | INTEGER NOT NULL           
 Name        | name         | TEXT NOT NULL             | TEXT NOT NULL              
 Fingerprint | fingerprint  | TEXT NOT NULL             | TEXT NOT NULL              
 Content     | content      | TEXT NOT NULL             | TEXT NOT NULL              
 Mode        | mode         | BIGINT NOT NULL DEFAULT 2 | INTEGER NOT NULL DEFAULT 2 
 Type        | type         | BIGINT NOT NULL DEFAULT 1 | INTEGER NOT NULL DEFAULT 1 
 CreatedUnix | created_unix | BIGINT                    | INTEGER                    
 UpdatedUnix | updated_unix | BIGINT                    | INTEGER                    

Primary keys: id
Indexes: 
	"idx_public_key_owner_id" (owner_id)
```

# Table "pull_request"

```
     Field      |      Column      | PostgreSQL  |        SQLite3        
----------------+------------------+-------------+-----------------------
 ID             | id               | BIGSERIAL   | INTEGER AUTOINCREMENT 
 Type           | type             | BIGINT      | INTEGER               
 Status         | status           | BIGINT      | INTEGER               
 IssueID        | issue_id         | BIGINT      | INTEGER               
 Index          | index            | BIGINT      | INTEGER               
 HeadRepoID     | head_repo_id     | BIGINT      | INTEGER               
 BaseRepoID     | base_repo_id     | BIGINT      | INTEGER               
 HeadUserName   | head_user_name   | TEXT        | TEXT                  
 HeadBranch     | head_branch      | TEXT        | TEXT                  
 BaseBranch     | base_branch      | TEXT        | TEXT                  
 MergeBase      | merge_base       | VARCHAR(40) | VARCHAR(40)           
 HasMerged      | has_merged       | BOOLEAN     | NUMERIC               
 MergedCommitID | merged_commit_id | VARCHAR(40) | VARCHAR(40)           
 MergerID       | merger_id        | BIGINT      | INTEGER               
 MergedUnix     | merged_unix      | BIGINT      | INTEGER               

Primary keys: id
Indexes: 
	"idx_pull_request_issue_id" (issue_id)
```

# Table "release"

```
    Field     |     Column     |           PostgreSQL           |            SQLite3             
--------------+----------------+--------------------------------+--------------------------------
 ID           | id             | BIGSERIAL                      | INTEGER AUTOINCREMENT          
 RepoID       | repo_id        | BIGINT                         | INTEGER                        
 PublisherID  | publisher_id   | BIGINT                         | INTEGER                        
 TagName      | tag_name       | TEXT NOT NULL                  | TEXT NOT NULL                  
 LowerTagName | lower_tag_name | TEXT NOT NULL                  | TEXT NOT NULL                  
 Target       | target         | TEXT NOT NULL                  | TEXT NOT NULL                  
 Title        | title          | TEXT                           | TEXT                           
 Sha1         | sha1           | VARCHAR(40)                    | VARCHAR(40)                    
 NumCommits   | num_commits    | BIGINT                         | INTEGER                        
 Note         | note           | TEXT                           | TEXT                           
 IsDraft      | is_draft       | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 IsPrerelease | is_prerelease  | BOOLEAN                        | NUMERIC                        
 CreatedUnix  | created_unix   | BIGINT                         | INTEGER                        

Primary keys: id
Indexes: 
	"idx_release_lower_tag_name" (lower_tag_name)
	"idx_release_publisher_id" (publisher_id)
	"idx_release_repo_id" (repo_id)
```

# Table "repository"

```
         Field         |          Column          |           PostgreSQL           |            SQLite3             
-----------------------+--------------------------+--------------------------------+--------------------------------
 ID                    | id                       | BIGSERIAL                      | INTEGER AUTOINCREMENT          
 OwnerID               | owner_id                 | BIGINT                         | INTEGER                        
 LowerName             | lower_name               | TEXT NOT NULL                  | TEXT NOT NULL                  
 Name                  | name                     | TEXT NOT NULL                  | TEXT NOT NULL                  
 Description           | description              | VARCHAR(512)                   | VARCHAR(512)                   
 Website               | website                  | TEXT                           | TEXT                           
 DefaultBranch         | default_branch           | TEXT                           | TEXT                           
 Size                  | size                     | BIGINT NOT NULL DEFAULT 0      | INTEGER NOT NULL DEFAULT 0     
 UseCustomAvatar       | use_custom_avatar        | BOOLEAN                        | NUMERIC                        
 NumWatches            | num_watches              | BIGINT                         | INTEGER                        
 NumStars              | num_stars                | BIGINT                         | INTEGER                        
 NumForks              | num_forks                | BIGINT                         | INTEGER                        
 NumIssues             | num_issues               | BIGINT                         | INTEGER                        
 NumClosedIssues       | num_closed_issues        | BIGINT                         | INTEGER                        
 NumPulls              | num_pulls                | BIGINT                         | INTEGER                        
 NumClosedPulls        | num_closed_pulls         | BIGINT                         | INTEGER                        
 NumMilestones         | num_milestones           | BIGINT NOT NULL DEFAULT 0      | INTEGER NOT NULL DEFAULT 0     
 NumClosedMilestones   | num_closed_milestones    | BIGINT NOT NULL DEFAULT 0      | INTEGER NOT NULL DEFAULT 0     
 IsPrivate             | is_private               | BOOLEAN                        | NUMERIC                        
 IsFork                | is_fork                  | BOOLEAN                        | NUMERIC                        
 IsUnlisted            | is_unlisted              | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 IsBare                | is_bare                  | BOOLEAN                        | NUMERIC                        
 IsMirror              | is_mirror                | BOOLEAN                        | NUMERIC                        
 EnableWiki            | enable_wiki              | BOOLEAN NOT NULL DEFAULT TRUE  | NUMERIC NOT NULL DEFAULT TRUE  
 AllowPublicWiki       | allow_public_wiki        | BOOLEAN                        | NUMERIC                        
 EnableExternalWiki    | enable_external_wiki     | BOOLEAN                        | NUMERIC                        
 ExternalWikiURL       | external_wiki_url        | TEXT                           | TEXT                           
 EnableIssues          | enable_issues            | BOOLEAN NOT NULL DEFAULT TRUE  | NUMERIC NOT NULL DEFAULT TRUE  
 AllowPublicIssues     | allow_public_issues      | BOOLEAN                        | NUMERIC                        
 EnableExternalTracker | enable_external_tracker  | BOOLEAN                        | NUMERIC                        
 ExternalTrackerURL    | external_tracker_url     | TEXT                           | TEXT                           
 ExternalTrackerFormat | external_tracker_format  | TEXT                           | TEXT                           
 ExternalTrackerStyle  | external_tracker_style   | TEXT                           | TEXT                           
 EnablePulls           | enable_pulls             | BOOLEAN NOT NULL DEFAULT TRUE  | NUMERIC NOT NULL DEFAULT TRUE  
 PullsIgnoreWhitespace | pulls_ignore_whitespace  | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 PullsAllowRebase      | pulls_allow_rebase       | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 PullsAllowManualMerge | pulls_allow_manual_merge | BOOLEAN NOT NULL DEFAULT FALSE | NUMERIC NOT NULL DEFAULT FALSE 
 PrimaryLanguage       | primary_language         | VARCHAR(32)                    | VARCHAR(32)                    
 Languages             | languages                | TEXT                           | TEXT                           
 ForkID                | fork_id                  | BIGINT                         | INTEGER                        
 CreatedUnix           | created_unix             | BIGINT                         | INTEGER                        
 UpdatedUnix           | updated_unix             | BIGINT                         | INTEGER                        

Primary keys: id
Indexes: 
	"idx_repository_lower_name" (lower_name)
	"idx_repository_name" (name)
	"repo_owner_name_unique" UNIQUE (owner_id, lower_name)
```

# Table "repo_topic"

```
  Field  |  Column  |   PostgreSQL    |     SQLite3      
---------+----------+-----------------+------------------
 TopicID | topic_id | BIGINT NOT NULL | INTEGER NOT NULL 
 RepoID  | repo_id  | BIGINT NOT NULL | INTEGER NOT NULL 

Primary keys: 
Indexes: 
	"repo_topic_unique" UNIQUE (topic_id, repo_id)
```

# Table "star"

```
 Field  | Column  |   PostgreSQL    |        SQLite3        
--------+---------+-----------------+-----------------------
 ID     | id      | BIGSERIAL       | INTEGER AUTOINCREMENT 
 UserID | uid     | BIGINT NOT NULL | INTEGER NOT NULL      
 RepoID | repo_id | BIGINT NOT NULL | INTEGER NOT NULL      

Primary keys: id
Indexes: 
	"star_user_repo_unique" UNIQUE (uid, repo_id)
```

# Table "team"

```
    Field    |   Column    |   PostgreSQL    |        SQLite3        
-------------+-------------+-----------------+-----------------------
 ID          | id          | BIGSERIAL       | INTEGER AUTOINCREMENT 
 OrgID       | org_id      | BIGINT          | INTEGER               
 LowerName   | lower_name  | TEXT NOT NULL   | TEXT NOT NULL         
 Name        | name        | TEXT NOT NULL   | TEXT NOT NULL         
 Description | description | VARCHAR(255)    | VARCHAR(255)          
 Authorize   | authorize   | BIGINT NOT NULL | INTEGER NOT NULL      
 NumRepos    | num_repos   | BIGINT          | INTEGER               
 NumMembers  | num_members | BIGINT          | INTEGER               

Primary keys: id
Indexes: 
	"idx_team_lower_name" (lower_name)
	"idx_team_org_id" (org_id)
```

# Table "team_repo"

```
 Field  | Column  | PostgreSQL |        SQLite3        
--------+---------+------------+-----------------------
 ID     | id      | BIGSERIAL  | INTEGER AUTOINCREMENT 
 OrgID  | org_id  | BIGINT     | INTEGER               
 TeamID | team_id | BIGINT     | INTEGER               
 RepoID | repo_id | BIGINT     | INTEGER               

Primary keys: id
Indexes: 
	"idx_team_repo_org_id" (org_id)
	"team_repo_team_repo_unique" UNIQUE (team_id, repo_id)
```

# Table "team_user"

```
 Field  | Column  | PostgreSQL |        SQLite3        
--------+---------+------------+-----------------------
 ID     | id      | BIGSERIAL  | INTEGER AUTOINCREMENT 
 OrgID  | org_id  | BIGINT     | INTEGER               
 TeamID | team_id | BIGINT     | INTEGER               
 UID    | uid     | BIGINT     | INTEGER               

Primary keys: id
Indexes: 
	"idx_team_user_org_id" (org_id)
	"team_user_team_user_unique" UNIQUE (team_id, uid)
```

# Table "topic"

```
    Field    |    Column    |      PostgreSQL      |        SQLite3        
-------------+--------------+----------------------+-----------------------
 ID          | id           | BIGSERIAL            | INTEGER AUTOINCREMENT 
 Name        | name         | TEXT NOT NULL UNIQUE | TEXT NOT NULL UNIQUE  
 NumRepos    | num_repos    | BIGINT               | INTEGER               
 CreatedUnix | created_unix | BIGINT               | INTEGER               
 UpdatedUnix | updated_unix | BIGINT               | INTEGER               

Primary keys: id
```

# Table "two_factor"

```
    Field    |    Column    |  PostgreSQL   |        SQLite3        
-------------+--------------+---------------+-----------------------
 ID          | id           | BIGSERIAL     | INTEGER AUTOINCREMENT 
 UserID      | user_id      | BIGINT UNIQUE | INTEGER UNIQUE        
 Secret      | secret       | TEXT          | TEXT                  
 CreatedUnix | created_unix | BIGINT        | INTEGER               

Primary keys: id
```

# Table "two_factor_recovery_code"

```
 Field  | Column  | PostgreSQL  |        SQLite3        
--------+---------+-------------+-----------------------
 ID     | id      | BIGSERIAL   | INTEGER AUTOINCREMENT 
 UserID | user_id | BIGINT      | INTEGER               
 Code   | code    | VARCHAR(11) | VARCHAR(11)           
 IsUsed | is_used | BOOLEAN     | NUMERIC               

Primary keys: id
Indexes: 
	"idx_two_factor_recovery_code_user_id" (user_id)
```

# Table "upload"

```
 Field | Column |      PostgreSQL      |        SQLite3        
-------+--------+----------------------+-----------------------
 ID    | id     | BIGSERIAL            | INTEGER AUTOINCREMENT 
 UUID  | uuid   | TEXT NOT NULL UNIQUE | TEXT NOT NULL UNIQUE  
 Name  | name   | TEXT NOT NULL        | TEXT NOT NULL         

Primary keys: id
```

# Table "user"

```
       Field        |        Column        |         PostgreSQL         |           SQLite3           
--------------------+----------------------+----------------------------+-----------------------------
 ID                 | id                   | BIGSERIAL                  | INTEGER AUTOINCREMENT       
 LowerName          | lower_name           | TEXT NOT NULL UNIQUE       | TEXT NOT NULL UNIQUE        
 Name               | name                 | TEXT NOT NULL              | TEXT NOT NULL               
 FullName           | full_name            | TEXT                       | TEXT                        
 Email              | email                | TEXT NOT NULL              | TEXT NOT NULL               
 Password           | passwd               | TEXT NOT NULL              | TEXT NOT NULL               
 LoginSource        | login_source         | BIGINT NOT NULL DEFAULT 0  | INTEGER NOT NULL DEFAULT 0  
 LoginName          | login_name           | TEXT                       | TEXT                        
 Type               | type                 | BIGINT                     | INTEGER                     
 Location           | location             | TEXT                       | TEXT                        
 Website            | website              | TEXT                       | TEXT                        
 Rands              | rands                | VARCHAR(10)                | VARCHAR(10)                 
 Salt               | salt                 | VARCHAR(10)                | VARCHAR(10)                 
 CreatedUnix        | created_unix         | BIGINT                     | INTEGER                     
 UpdatedUnix        | updated_unix         | BIGINT                     | INTEGER                     
 LastRepoVisibility | last_repo_visibility | BOOLEAN                    | NUMERIC                     
 MaxRepoCreation    | max_repo_creation    | BIGINT NOT NULL DEFAULT -1 | INTEGER NOT NULL DEFAULT -1 
 IsActive           | is_active            | BOOLEAN                    | NUMERIC                     
 IsAdmin            | is_admin             | BOOLEAN                    | NUMERIC                     
 AllowGitHook       | allow_git_hook       | BOOLEAN                    | NUMERIC                     
 AllowImportLocal   | allow_import_local   | BOOLEAN                    | NUMERIC                     
 ProhibitLogin      | prohibit_login       | BOOLEAN                    | NUMERIC                     
 Avatar             | avatar               | VARCHAR(2048) NOT NULL     | VARCHAR(2048) NOT NULL      
 AvatarEmail        | avatar_email         | TEXT NOT NULL              | TEXT NOT NULL               
 UseCustomAvatar    | use_custom_avatar    | BOOLEAN                    | NUMERIC                     
 NumFollowers       | num_followers        | BIGINT                     | INTEGER                     
 NumFollowing       | num_following        | BIGINT NOT NULL DEFAULT 0  | INTEGER NOT NULL DEFAULT 0  
 NumStars           | num_stars            | BIGINT                     | INTEGER                     
 NumRepos           | num_repos            | BIGINT                     | INTEGER                     
 Description        | description          | TEXT                       | TEXT                        
 NumTeams           | num_teams            | BIGINT                     | INTEGER                     
 NumMembers         | num_members          | BIGINT                     | INTEGER                     
 DisableNewsFeed    | disable_news_feed    | BOOLEAN                    | NUMERIC                     

Primary keys: id
```

# Table "version"

```
  Field  | Column  |   PostgreSQL    |        SQLite3        
---------+---------+-----------------+-----------------------
 ID      | id      | BIGSERIAL       | INTEGER AUTOINCREMENT 
 Version | version | BIGINT NOT NULL | INTEGER NOT NULL      

Primary keys: id
```

# Table "watch"

```
 Field  | Column  |   PostgreSQL    |        SQLite3        
--------+---------+-----------------+-----------------------
 ID     | id      | BIGSERIAL       | INTEGER AUTOINCREMENT 
 UserID | user_id | BIGINT NOT NULL | INTEGER NOT NULL      
 RepoID | repo_id | BIGINT NOT NULL | INTEGER NOT NULL      

Primary keys: id
Indexes: 
	"watch_user_repo_unique" UNIQUE (user_id, repo_id)
```

