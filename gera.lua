for i=1,100 do
	f = io.open("in/t" ..i , "w")
	
	f:write(string.rep("aaaaaaa\n",1024*1024*10))
	
	f:close()
end
