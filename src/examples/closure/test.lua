--
-- Created by IntelliJ IDEA.
-- User: RonanLuo
-- Date: 2020/10/14
-- Time: 12:13
-- To change this template use File | Settings | File Templates.
--

local function max(...)
    local args = { ... }
    local val, idx
    for i = 1, #args do
        if val == nil or args[i] > val then
            val, idx = args[i], i
        end
    end
    print("max", val, idx)
    return val, idx
end

local function assert(v)
    print("assert", v)
    if not v then print("assert failed") end
end

local v1 = max(3, 9, 7, 128, 35)
assert(v1 == 128)
local v2, i2 = max(3, 9, 7, 128, 35)
assert(v2 == 128 and i2 == 4)
local v3, i3 = max(max(3, 9, 7, 128, 35))
assert(v3 == 128 and i3 == 1)
local t = { max(3, 9, 7, 128, 35) }
assert(t[1] == 128 and t[2] == 4)
print("hello world", v1, v2, v3, t)
