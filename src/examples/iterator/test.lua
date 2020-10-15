--
-- Created by IntelliJ IDEA.
-- User: RonanLuo
-- Date: 2020/10/15
-- Time: 16:09
-- To change this template use File | Settings | File Templates.
--
t = {a = 1, b = 2, c = 3}
for k, v in pairs(t) do
    print(k, v)
end

t = {"a", "b", "c"}
for k, v in ipairs(t) do
    print(k, v)
end


