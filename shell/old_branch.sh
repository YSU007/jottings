#!/bin/bash

# 设置参数
repository_path=$(git rev-parse --show-toplevel)  # 获取当前目录的Git仓库根路径
repo_name=$(basename "$repository_path")  # 获取仓库名称
report_file="/tmp/report_${repo_name}_$(date -Iseconds).csv"  # 输出报告的文件名

# 进入仓库目录
cd "$repository_path" || exit 1

# 修改为你希望的阈值天数
days_threshold=100
if [ -n "$1" ]; then
    days_threshold=$1
fi

# 获取远程分支列表
remote_branches=$(git for-each-ref --format='%(refname:short)' refs/remotes/origin/)

# 获取当前日期的时间戳
current_timestamp=$(date +%s)

# 创建或清空报告文件
echo "Branch,Last Commit Date,Days Since Last Commit,Last Commit Hash,Last Commit Message,Last Commit Auther,JWS2 Link" > "$report_file"

# 循环遍历远程分支
for remote_branch in $remote_branches; do
    branch_name=${remote_branch#origin/}  # 去除前缀
    # 如果分支名称为 "HEAD" 或以 "live/" "release/" 开头，跳过这个分支
    if [[ $branch_name == "HEAD" || $branch_name == live/* || $branch_name == release/* || $branch_name == develop ]]; then
        continue
    fi
    # 获取远程分支的最后一次提交时间戳
    last_commit_timestamp=$(git log -n 1 --format="%at" "$remote_branch")

    # 计算时间差
    time_diff=$((current_timestamp - last_commit_timestamp))
    days_diff=$((time_diff / 86400))  # 一天有 86400 秒

    # 检查是否超过阈值天数
    if [ $days_diff -gt $days_threshold ]; then
        # 获取最近一次提交的日期
        last_commit_date=$(git log -n 1 --format="%cd" --date=short "$remote_branch")
        last_commit_hash=$(git log -n 1 --format="%H" "$remote_branch")
        latest_commit_msg=$(git log -n 1 --format="%s" "$remote_branch")

        # 获取最近一次提交的作者
        last_commit_auther=$(git log -n 1 --format="%an" "$remote_branch")

        # 判断分支名或提交信息是否包含 "JWS2-" 类似的关键字不区分大小写，然后从包含关键字的分支名或提交中提取JWS2编号，优先取分支名中的
        jws2_number=""
        if [[ $branch_name == *"JWS2-"* ]]; then
            # 使用 awk 提取 JWS2 编号
            jws2_number=$(echo "$branch_name" | awk -F'JWS2-' '{print $2}' | awk '{print $1}' | sed 's/[^0-9].*//')
        elif [[ $latest_commit_msg == *"JWS2-"* ]]; then
            # 使用 awk 提取 JWS2 编号
            jws2_number=$(echo "$latest_commit_msg" | awk -F'JWS2-' '{print $2}' | awk '{print $1}' | sed 's/[^0-9].*//')
        fi
        jws2_link=""
        if [ -n "$jws2_number" ]; then
            jws2_link="http://jira-office.taiplat.com:18080/browse/JWS2-$jws2_number"
        fi

        # 打印部分信息到控制台,方便调试
        echo "Scan branch: $branch_name $jws2_link"

        # 输出到报告文件
        printf "%s,%s,%d,\"%s\",\"%s\",\"%s\",\"%s\"\n" \
            "$branch_name" "$last_commit_date" "$days_diff" "$last_commit_hash" \
            "$latest_commit_msg" "$last_commit_auther" "$jws2_link" >> "$report_file"
    fi
done