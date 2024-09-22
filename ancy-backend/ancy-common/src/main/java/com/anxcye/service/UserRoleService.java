package com.anxcye.service;

import com.anxcye.domain.entity.UserRole;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_user_role(用户和角色关联表)】的数据库操作Service
* @createDate 2024-09-20 14:19:43
*/
public interface UserRoleService extends IService<UserRole> {


    boolean addByUserId(Long userId, List<Long> roleIds);

    boolean deleteByUserId(Long userId);

    boolean deleteByRoleId(Long roleId);

    List<Long> selectRoleIdsByUserId(Long id);

    boolean updateByUserId(Long id, List<Long> roleIds);
}
