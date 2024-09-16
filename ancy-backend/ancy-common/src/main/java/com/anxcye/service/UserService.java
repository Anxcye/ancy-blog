package com.anxcye.service;

import com.anxcye.domain.entity.User;
import com.anxcye.domain.vo.BlogUserVo;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【sys_user(用户表)】的数据库操作Service
* @createDate 2024-09-12 13:57:34
*/
public interface UserService extends IService<User> {

    BlogUserVo login(User user);
}
