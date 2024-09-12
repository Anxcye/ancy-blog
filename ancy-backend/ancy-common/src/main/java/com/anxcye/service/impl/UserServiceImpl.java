package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.User;
import com.anxcye.service.UserService;
import com.anxcye.mapper.UserMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【sys_user(用户表)】的数据库操作Service实现
* @createDate 2024-09-12 13:57:34
*/
@Service
public class UserServiceImpl extends ServiceImpl<UserMapper, User>
    implements UserService{

}




