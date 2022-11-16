
count = 0

--  count = count+1
 if count == 1 then
   print("ok!")
else
   print(count)
end


function main()
  print("Hello from Lua!")
end

-- function fib(n)
-- 	if n < 2 then
-- 		return 1
-- 	end
-- 	return fib(n-1) + fib(n-2)
-- end
-- local time = require("time")

-- local tick = time:ticker(1000)
-- local count = 0

-- time:after(5000,function()
--     print("after 5s")
-- end)

-- repeat
--     tick:wait()
--     count = count + 1
--     print(count)
-- until(count > 10)

-- tick:close()