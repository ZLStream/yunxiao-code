#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');

function removeExisting(target) {
  try {
    const stat = fs.lstatSync(target);
    if (stat.isSymbolicLink()) {
      fs.unlinkSync(target);
    } else if (stat.isDirectory()) {
      fs.rmSync(target, { recursive: true });
    } else {
      fs.unlinkSync(target);
    }
    return true;
  } catch (e) {
    if (e.code === 'ENOENT') return true;
    return false;
  }
}

function deployPlugin() {
  const claudeDir = path.join(os.homedir(), '.claude');
  const commandsDir = path.join(claudeDir, 'commands');
  const skillsDir = path.join(claudeDir, 'skills');

  // 获取 npm 包路径
  const packageDir = path.dirname(require.main.filename);
  const pluginDir = path.join(packageDir, '..', 'claude-plugin');

  // 确保目标目录存在
  if (!fs.existsSync(commandsDir)) {
    fs.mkdirSync(commandsDir, { recursive: true });
  }
  if (!fs.existsSync(skillsDir)) {
    fs.mkdirSync(skillsDir, { recursive: true });
  }

  const commandsTarget = path.join(commandsDir, 'yx-commands');
  const skillsTarget = path.join(skillsDir, 'yx-commands');
  const commandsSource = path.join(pluginDir, 'commands');
  const skillsSource = path.join(pluginDir, 'skills');

  // 创建/更新 commands 软链接
  if (fs.existsSync(commandsSource)) {
    if (removeExisting(commandsTarget)) {
      fs.symlinkSync(commandsSource, commandsTarget);
      console.log('✅ 已链接 commands 到 ~/.claude/commands/yx-commands/');
    }
  }

  // 创建/更新 skills 软链接
  if (fs.existsSync(skillsSource)) {
    if (removeExisting(skillsTarget)) {
      fs.symlinkSync(skillsSource, skillsTarget);
      console.log('✅ 已链接 skills 到 ~/.claude/skills/yx-commands/');
    }
  }

  console.log('\n可用命令: /yx-commands:commit, /yx-commands:push, /yx-commands:mr, /yx-commands:review\n');
}

console.log('\n✅ yx-code CLI installed successfully!\n');
console.log('Usage:');
console.log('  yx-code init     - 初始化云效配置');
console.log('  yx-code clone    - 克隆仓库');
console.log('  yx-code commit   - 提交代码');
console.log('  yx-code push     - 推送代码');
console.log('  yx-code mr       - 创建合并请求');
console.log('  yx-code diff     - 查看代码差异\n');

try {
  deployPlugin();
} catch (err) {
  console.log('⚠️  部署失败:', err.message, '\n');
}