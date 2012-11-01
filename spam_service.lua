require "math"

function die( msg )
	print( msg )
	os.exit(1)
end

function delay_s(delay)
    delay = delay or 1
    local time_to = os.time() + delay
    while os.time() < time_to do end
end

function exec_or_die( cmd )
	if os.execute( cmd ) ~= 0 then
		die( "Failed to execute \"" .. cmd .. "\"" )
	end
end

local tries = 0
while ( tries < 6000 ) do
	exec_or_die( 'call curl "http://localhost:2200/reading/create?key=FLOW_01&val=' .. math.sin(tries%60)/60*13 + math.random()*2 .. '"' )
	exec_or_die( 'call curl "http://localhost:2200/reading/create?key=PRES_01&val=' .. math.sin(tries%60)/60*13 + math.random()*30 .. '"' )
	delay_s( 1 )
	tries = tries + 1
end