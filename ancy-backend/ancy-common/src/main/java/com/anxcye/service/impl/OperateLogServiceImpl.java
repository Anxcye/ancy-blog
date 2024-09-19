package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.OperateLog;
import com.anxcye.service.OperateLogService;
import com.anxcye.mapper.OperateLogMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【sys_operate_log(操作日志表)】的数据库操作Service实现
* @createDate 2024-09-19 16:16:28
*/
@Service
public class OperateLogServiceImpl extends ServiceImpl<OperateLogMapper, OperateLog>
    implements OperateLogService{

}




