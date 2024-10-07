package com.anxcye.service.impl;

import com.anxcye.annotation.Log;
import com.anxcye.domain.entity.UserRole;
import com.anxcye.mapper.UserRoleMapper;
import com.anxcye.service.UserRoleService;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

/**
* @author axy
* @description 针对表【sys_user_role(用户和角色关联表)】的数据库操作Service实现
* @createDate 2024-09-20 14:19:43
*/
@Service
public class UserRoleServiceImpl extends ServiceImpl<UserRoleMapper, UserRole>
    implements UserRoleService{

    @Override
    @Transactional
    public boolean addByUserId(Long userId, List<Long> roleIds) {
        if (roleIds == null || roleIds.isEmpty()) {
            return true;
        }
        List<UserRole> userRoles = roleIds.stream()
                .map(roleId -> new UserRole(userId, roleId))
                .collect(Collectors.toList());
        return saveBatch(userRoles);
    }

    @Override
    public boolean deleteByUserId(Long userId) {
        LambdaQueryWrapper<UserRole> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserRole::getUserId, userId);
        return remove(wrapper);
    }

    @Override
    public boolean deleteByRoleId(Long roleId) {
        LambdaQueryWrapper<UserRole> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserRole::getRoleId, roleId);
        return remove(wrapper);
    }

    @Override
    public List<Long> selectRoleIdsByUserId(Long userId) {
        LambdaQueryWrapper<UserRole> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(UserRole::getUserId, userId);
        List<UserRole> userRoles = list(wrapper);
        return userRoles.stream()
                .map(UserRole::getRoleId)
                .collect(Collectors.toList());
    }

    @Log
    @Override
    @Transactional
    public boolean updateByUserId(Long id, List<Long> roleIds) {
        deleteByUserId(id);
        addByUserId(id, roleIds);
        return true;
    }
}




