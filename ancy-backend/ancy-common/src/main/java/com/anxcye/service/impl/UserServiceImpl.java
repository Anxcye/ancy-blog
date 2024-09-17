package com.anxcye.service.impl;

import com.anxcye.constants.RedisConstant;
import com.anxcye.domain.entity.LoginUser;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.vo.BlogUserVo;
import com.anxcye.domain.vo.UserInfoVo;
import com.anxcye.exception.SystemException;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.JwtUtil;
import com.anxcye.utils.RedisCache;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.User;
import com.anxcye.service.UserService;
import com.anxcye.mapper.UserMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;

import java.util.Objects;

/**
 * @author axy
 * @description 针对表【sys_user(用户表)】的数据库操作Service实现
 * @createDate 2024-09-12 13:57:34
 */
@Service
public class UserServiceImpl extends ServiceImpl<UserMapper, User>
        implements UserService {

    @Autowired
    private AuthenticationManager authenticationManager;

    @Autowired
    private RedisCache redisCache;

    @Override
    public BlogUserVo login(User user) {
        UsernamePasswordAuthenticationToken authenticationToken =
                new UsernamePasswordAuthenticationToken(user.getUserName(), user.getPassword());
        Authentication authenticate = authenticationManager.authenticate(authenticationToken);

        if (Objects.isNull(authenticate)) {
            throw new SystemException(AppHttpCodeEnum.LOGIN_ERROR);
        }

        LoginUser loginUser = (LoginUser) authenticate.getPrincipal();
        String id = loginUser.getUser().getId().toString();
        String jwt = JwtUtil.createJWT(id);

        redisCache.setCacheObject(RedisConstant.BLOG_TOKEN_PREFIX + id, loginUser);

        UserInfoVo userInfoVo = BeanCopyUtils.copyBean(loginUser.getUser(), UserInfoVo.class);

        return new BlogUserVo(jwt, userInfoVo);

    }

    @Override
    public void logout() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        LoginUser loginUser = (LoginUser) authentication.getPrincipal();

        Long id = loginUser.getUser().getId();
        redisCache.deleteObject(RedisConstant.BLOG_TOKEN_PREFIX + id);
    }
}




