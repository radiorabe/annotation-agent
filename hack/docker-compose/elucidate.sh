#!/bin/bash

AUTH_OPTS=" -Dauth.enabled=$AUTH_ENABLED -Dauth.token.verifierKey='$AUTH_VERIFIER_KEY' -Dauth.token.verifierType=$AUTH_VERIFIER_TYPE -Dauth.anonReadAccess=true "

export CATALINA_OPTS=" -Dlog4j.config.location=$LOG4J_CONFIG_LOCATION -Ddb.url=$DATABASE_URL -Ddb.user=$DATABASE_USER -Ddb.password=$DATABASE_PASSWORD -Dbase.path=$BASE_PATH $AUTH_OPTS "

exec /usr/local/tomcat/bin/catalina.sh run $@
