package com.anxcye.service.impl;

import com.anxcye.domain.entity.Role;
import com.anxcye.mapper.RoleMapper;
import com.anxcye.service.RoleService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author axy
 * @description 针对表【sys_role(角色信息表)】的数据库操作Service实现
 * @createDate 2024-09-20 13:25:43
 */
@Service
public class RoleServiceImpl extends ServiceImpl<RoleMapper, Role>
        implements RoleService {

    @Override
    public List<String> getRoleByUserId(Long userId) {
        return getBaseMapper().getRolesByUserId(userId).stream()
                .map(Role::getRoleKey).toList();

    }
}




