--
-- Created by IntelliJ IDEA.
-- User: RonanLuo
-- Date: 2020/10/15
-- Time: 9:35
-- To change this template use File | Settings | File Templates.
--

function newCounter ()
    local count = 0
    return function () -- anonymous function
        count = count + 1
        return count
    end
end

c1 = newCounter()
print(c1()) --> 1
print(c1()) --> 2

c2 = newCounter()
print(c2()) --> 1
print(c1()) --> 3
print(c2()) --> 2

