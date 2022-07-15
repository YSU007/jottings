#!/bin/bash

Checkout(){
  Branch=$1
  TargetDir=$2
  # 切换到项目目录
  cd "$TargetDir" || exit
  # 执行git命令
  # 将当前的分支stash
  if [ -n "$(git status -s)" ];then
      git stash save -m "$(date)"
  fi
  # 切换到Branch分支
  git checkout "$Branch"
  # 更新
  git pull -ff
}

HardReset(){
  Checkout "$1" "$2"
  git reset --hard "origin/$1"
}

NewBranch() {
    BaseBranch=$1
    TargetBranch=$2
    TargetDir=$3
    # 保存当前目录
    echo "Start to create branch"
    echo "$TargetBranch" "$TargetDir"
    # 切分支
    Checkout "$BaseBranch" "$TargetDir"
    # 更新
    git pull
    # 新建的分支
    git branch "$TargetBranch"
    # 切换到新建的分支
    git checkout "$TargetBranch"
    # 显示创建的分支，看是否已经切换成功
    git branch -vv
    # stage新建分支信息
    git status
    git add .
    git commit -m "branch create $TargetBranch"
    # 将分支上传到服务器
    git push origin "$TargetBranch"
    # 设置依赖分支
    git branch --set-upstream-to="origin/$TargetBranch"
    # 显示全部服务器分支
    git branch -a | grep -i "$TargetBranch"
    # 删除缓存
    echo "Create branch $TargetBranch Success";
}

DelBranch() {
    BaseBranch=$1
    TargetBranch=$2
    TargetDir=$3
    echo "Start to del branch"
    echo "$TargetBranch" "$TargetDir"
    # 切分支
    Checkout "$BaseBranch" "$TargetDir"
    # 更新
    git pull
    # 删本地
    git branch --delete --force "$TargetBranch"
    # 删远程
    git push origin --delete "refs/heads/$TargetBranch"
    echo "Del  $TargetBranch Success";
}