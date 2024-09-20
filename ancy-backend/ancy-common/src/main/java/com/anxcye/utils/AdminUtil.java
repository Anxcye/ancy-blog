package com.anxcye.utils;

import com.anxcye.constants.SystemConstants;

public class AdminUtil {
    public static boolean isSuperAdmin(Long id){
        return id == SystemConstants.SUPER_ADMIN_ID;
    }
}
