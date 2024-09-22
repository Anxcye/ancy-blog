package com.anxcye.utils;

import com.anxcye.domain.entity.LoginUser;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

public class SecurityUtil {

    /**
     * 获取用户
     **/
    public static LoginUser getLoginUser() {
        return (LoginUser) getAuthentication().getPrincipal();
    }

    /**
     * 获取Authentication
     */
    public static Authentication getAuthentication() {
        return SecurityContextHolder.getContext().getAuthentication();
    }

    public static Boolean isAdmin() {
        try {
            Long id = getLoginUser().getUser().getId();
            return id != null && 1L == id;
        } catch (Exception e) {
            return false;
        }
    }

    public static Long getUserId() {
        return getLoginUser().getUser().getId();
    }
}
