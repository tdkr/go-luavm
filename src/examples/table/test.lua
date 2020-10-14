--
-- Created by IntelliJ IDEA.
-- User: RonanLuo
-- Date: 2020/10/14
-- Time: 9:52
-- To change this template use File | Settings | File Templates.
--

local t = { "a", "b", "c" }
t[2] = "B"
t["foo"] = "Bar"
local s = t[3] .. t[2] .. t[1] .. t["foo"] .. #t
