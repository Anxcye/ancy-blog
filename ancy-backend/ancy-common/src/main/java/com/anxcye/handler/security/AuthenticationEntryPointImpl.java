package com.anxcye.handler.security;

import com.alibaba.fastjson.JSON;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.utils.WebUtils;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.InsufficientAuthenticationException;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class AuthenticationEntryPointImpl implements AuthenticationEntryPoint {
    @Override
    public void commence(HttpServletRequest request, HttpServletResponse response, AuthenticationException authException) throws IOException, ServletException {
        authException.printStackTrace();

        ResponseResult result = null;
        if (authException instanceof InsufficientAuthenticationException) {
            result = ResponseResult.error(AppHttpCodeEnum.NEED_LOGIN);
        } else if (authException instanceof BadCredentialsException) {
            result = ResponseResult.error(AppHttpCodeEnum.LOGIN_ERROR);
        } else {
            result = ResponseResult.error(AppHttpCodeEnum.AUTH_ERROR);
        }
        WebUtils.renderString(response, JSON.toJSONString(result));
    }
}
