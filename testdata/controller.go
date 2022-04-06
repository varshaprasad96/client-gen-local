package controller

// +custom:generate:types=memcached,groups=batch.io,resources=cronjobs,verbs=get;watch;create
// +custom:generate:types=something,groups=batch.io,resources=cronjobs/status,verbs=get;update;patch
// +custom:generate:types=test,groups=art,resources=jobs,verbs=get
