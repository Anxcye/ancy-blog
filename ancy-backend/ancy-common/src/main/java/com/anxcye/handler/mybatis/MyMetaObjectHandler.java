package com.anxcye.handler.mybatis;


import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.handlers.MetaObjectHandler;
import org.apache.ibatis.reflection.MetaObject;
import org.springframework.stereotype.Component;

import java.util.Date;

@Component
public class MyMetaObjectHandler implements MetaObjectHandler {
    @Override
    public void insertFill(MetaObject metaObject) {
        try {
            Long userId = SecurityUtil.getUserId();
            this.setFieldValByName("createTime", new Date(), metaObject);
            this.setFieldValByName("createBy", userId, metaObject);
            this.setFieldValByName("updateTime", new Date(), metaObject);
            this.setFieldValByName("updateBy", userId, metaObject);
        } catch (Exception e) {

        }
    }

    @Override
    public void updateFill(MetaObject metaObject) {
        try {
            Long userId = SecurityUtil.getUserId();
            this.setFieldValByName("updateTime", new Date(), metaObject);
            this.setFieldValByName("updateBy", userId, metaObject);
        } catch (Exception e) {

        }
    }
}