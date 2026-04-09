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
  const skillsDir = path.join(claudeDir, 'skills');

  // 获取 npm 包路径
  const packageDir = path.dirname(require.main.filename);
  const pluginDir = path.join(packageDir, '..', 'claude-plugin');
  const skillsSource = path.join(pluginDir, 'skills');

  // 确保目标目录存在
  if (!fs.existsSync(skillsDir)) {
    fs.mkdirSync(skillsDir, { recursive: true });
  }

  // 清理旧的 commands 软链接（如果存在）
  const oldCommandsTarget = path.join(claudeDir, 'commands', 'yx-commands');
  removeExisting(oldCommandsTarget);

  // 清理旧版本的单一 skills 软链接（指向整个 skills 目录的）
  const oldSkillsTarget = path.join(skillsDir, 'yx-commands');
  try {
    const linkTarget = fs.readlinkSync(oldSkillsTarget);
    if (linkTarget === skillsSource || linkTarget.endsWith('/claude-plugin/skills')) {
      removeExisting(oldSkillsTarget);
    }
  } catch (e) { /* not a symlink or doesn't exist */ }

  if (!fs.existsSync(skillsSource)) {
    console.log('⚠️  skills 目录不存在:', skillsSource);
    return;
  }

  // 为每个 skill 子目录单独创建符号链接
  const entries = fs.readdirSync(skillsSource, { withFileTypes: true });
  const linked = [];
  for (const entry of entries) {
    if (!entry.isDirectory()) continue;
    const source = path.join(skillsSource, entry.name);
    const target = path.join(skillsDir, entry.name);
    if (removeExisting(target)) {
      fs.symlinkSync(source, target);
      linked.push(entry.name);
    }
  }

  if (linked.length > 0) {
    console.log('✅ 已链接 ' + linked.length + ' 个 skills 到 ~/.claude/skills/');
    linked.forEach(name => console.log('   - ' + name));
  }

  console.log('\n可用命令: /yx-commands-commit, /yx-commands-push, /yx-commands-mr, /yx-commands-review, /yx-commands-commit-push-mr, /yx-commands-workitems, /yx-commands-summary\n');
}

console.log('\n✅ yx-code CLI installed successfully!\n');
console.log('Usage:');
console.log('  yx-code init     - 初始化云效配置');
console.log('  yx-code clone    - 克隆仓库');
console.log('  yx-code commit   - 提交代码');
console.log('  yx-code push     - 推送代码');
console.log('  yx-code mr       - 创建合并请求');
console.log('  yx-code diff     - 查看代码差异');
console.log('  yx-code workitems - 列出工作项');
console.log('  yx-code summary  - 生成工作总结\n');

try {
  deployPlugin();
} catch (err) {
  console.log('⚠️  部署失败:', err.message, '\n');
}