
 local count
-- 如果 count 为空则等于 0
if count == nil then
	count = 0
end


-- 全局变量 0 如果不存在就创建等于 1
 if not _G["count"] then
	 _G["count"] = 1
 else
	 _G["count"] = _G["count"] + 1
 end

 if not c then
	c = 1
 else
	c = c + 1
 end


 count = count+1
 if count == 1 then
   print("ok!",count,c)
else
   print('open : ', count)
end


function main()
  print("Hello from Lua!")
end

function fib2(n)
  print("Hello from fib2!")
return fib(n)
end
function fib(n)
	if n < 2 then
		return 1
	end
	return fib(n-1) + fib(n-2)
end
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