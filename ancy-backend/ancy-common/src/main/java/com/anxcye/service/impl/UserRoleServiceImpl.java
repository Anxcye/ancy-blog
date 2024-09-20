package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.UserRole;
import com.anxcye.service.UserRoleService;
import com.anxcye.mapper.UserRoleMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【sys_user_role(用户和角色关联表)】的数据库操作Service实现
* @createDate 2024-09-20 14:19:43
*/
@Service
public class UserRoleServiceImpl extends ServiceImpl<UserRoleMapper, UserRole>
    implements UserRoleService{

}




