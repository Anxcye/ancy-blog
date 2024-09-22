package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.LoginUser;
import com.anxcye.domain.entity.User;
import com.anxcye.mapper.UserMapper;
import com.anxcye.service.MenuService;
import com.anxcye.service.RoleService;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;

@Service
public class UserDetailServiceImpl implements UserDetailsService {

    @Autowired
    private UserMapper userMapper;

    @Autowired
    private MenuService menuService;

    @Autowired
    private RoleService roleService;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        LambdaQueryWrapper<User> userLambdaQueryWrapper = new LambdaQueryWrapper<>();
        userLambdaQueryWrapper.eq(User::getUserName, username);

        User user = userMapper.selectOne(userLambdaQueryWrapper);

        if (Objects.isNull(user)) {
            throw new RuntimeException("用户不存在");
        }

        // TODO: auth
        if (Objects.equals(user.getType(), SystemConstants.USER_ADMIN)) {
            List<String> permissions = menuService.getPermissionsByUserId(user.getId());
            List<String> roles = roleService.getRoleByUserId(user.getId());

            return new LoginUser(user, permissions, roles);
        }

        return new LoginUser(user);
    }


}
