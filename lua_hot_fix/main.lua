function sleep(n)
    local time = os.clock();
    while true do
        if os.clock() - time > n then
            return
        end
    end
end

local HotFix = require("hotfix")
require("fix")

while true do
    fixfunc()

    sleep(2)
    -- 开始热更
    HotFix:UpdateModule("fix")
end