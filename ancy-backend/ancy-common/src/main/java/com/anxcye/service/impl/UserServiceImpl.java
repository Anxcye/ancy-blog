package com.anxcye.service.impl;

import com.anxcye.constants.RedisConstant;
import com.anxcye.domain.dto.LoginDto;
import com.anxcye.domain.dto.RegisterDto;
import com.anxcye.domain.dto.UserDto;
import com.anxcye.domain.entity.LoginUser;
import com.anxcye.domain.entity.User;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.vo.AdminUserVo;
import com.anxcye.domain.vo.BlogUserVo;
import com.anxcye.domain.vo.UserInfoVo;
import com.anxcye.exception.SystemException;
import com.anxcye.mapper.UserMapper;
import com.anxcye.service.MenuService;
import com.anxcye.service.RoleService;
import com.anxcye.service.UserService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.JwtUtil;
import com.anxcye.utils.RedisCache;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.toolkit.support.SFunction;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;

/**
 * @author axy
 * @description 针对表【sys_user(用户表)】的数据库操作Service实现
 * @createDate 2024-09-12 13:57:34
 */
@Service
public class UserServiceImpl extends ServiceImpl<UserMapper, User> implements UserService {

    @Autowired
    private AuthenticationManager authenticationManager;

    @Autowired
    private RedisCache redisCache;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private MenuService menuService;

    @Autowired
    private RoleService roleService;

    private LoginUser getLoginUserByLoginDto(LoginDto loginDto) {
        if (Objects.isNull(loginDto.getUserName())) {
            loginDto.setUserName(getOne(new LambdaQueryWrapper<User>().eq(User::getEmail, loginDto.getUserName())).getUserName());
        }

        UsernamePasswordAuthenticationToken authenticationToken = new UsernamePasswordAuthenticationToken(loginDto.getUserName(), loginDto.getPassword());
        Authentication authenticate = authenticationManager.authenticate(authenticationToken);

        if (Objects.isNull(authenticate)) {
            throw new SystemException(AppHttpCodeEnum.LOGIN_ERROR);
        }

        return (LoginUser) authenticate.getPrincipal();
    }

    @Override
    public BlogUserVo login(LoginDto loginDto) {

        LoginUser loginUser = getLoginUserByLoginDto(loginDto);

        String id = loginUser.getUser().getId().toString();

        redisCache.setCacheObject(RedisConstant.BLOG_TOKEN_PREFIX + id, loginUser);

        String jwt = JwtUtil.createJWT(id);
        UserInfoVo userInfoVo = BeanCopyUtils.copyBean(loginUser.getUser(), UserInfoVo.class);

        return new BlogUserVo(jwt, userInfoVo);

    }

    @Override
    public void logout() {
        LoginUser loginUser = SecurityUtil.getLoginUser();

        Long id = loginUser.getUser().getId();
        redisCache.deleteObject(RedisConstant.BLOG_TOKEN_PREFIX + id);
    }

    @Override
    public UserInfoVo getUserInfo() {
        Long userId = SecurityUtil.getUserId();
        User user = getById(userId);
        return BeanCopyUtils.copyBean(user, UserInfoVo.class);
    }

    @Override
    public UserInfoVo updateUserInfo(UserDto userDto) {
        if (!Objects.equals(SecurityUtil.getUserId(), userDto.getId())) {
            throw new SystemException(AppHttpCodeEnum.NO_OPERATOR_AUTH);
        }
        User user = BeanCopyUtils.copyBean(userDto, User.class);
        updateById(user);
        return getUserInfo();
    }

    private boolean userInfoExists(SFunction<User, String> getter, String value) {
        return getOne(new LambdaQueryWrapper<User>().eq(getter, value)) != null;
    }

    @Override
    public BlogUserVo register(RegisterDto userDto) {
        if (Objects.isNull(userDto.getUserName()) || Objects.isNull(userDto.getPassword()) || Objects.isNull(userDto.getEmail()) || Objects.isNull(userDto.getNickName())) {
            throw new SystemException(AppHttpCodeEnum.USERINFO_NOT_NULL);
        }
        if (userInfoExists(User::getUserName, userDto.getUserName())) {
            throw new SystemException(AppHttpCodeEnum.USERNAME_EXIST);
        }
        if (userInfoExists(User::getEmail, userDto.getEmail())) {
            throw new SystemException(AppHttpCodeEnum.EMAIL_EXIST);
        }

        User user = BeanCopyUtils.copyBean(userDto, User.class);
        String encodedPassword = passwordEncoder.encode(userDto.getPassword());
        user.setPassword(encodedPassword);

        save(user);
        LoginDto loginDto = BeanCopyUtils.copyBean(userDto, LoginDto.class);
        return login(loginDto);
    }

    @Override
    public AdminUserVo adminLogin(LoginDto loginDto) {

        LoginUser loginUser = getLoginUserByLoginDto(loginDto);
        Long id = loginUser.getUser().getId();
        redisCache.setCacheObject(RedisConstant.ADMIN_TOKEN_PREFIX + id, loginUser);

        String jwt = JwtUtil.createJWT(id.toString());
        UserInfoVo userInfoVo = BeanCopyUtils.copyBean(loginUser.getUser(), UserInfoVo.class);
        List<String> permissions = menuService.getPermissionsByUserId(id);
        List<String> role = roleService.getRoleByUserId(id);


        return new AdminUserVo(jwt, permissions, role, userInfoVo);
    }
}
