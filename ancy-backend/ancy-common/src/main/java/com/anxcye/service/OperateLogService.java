package com.anxcye.service;

import com.anxcye.domain.dto.TimelinePageDto;
import com.anxcye.domain.entity.OperateLog;
import com.anxcye.domain.result.PageResult;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【sys_operate_log(操作日志表)】的数据库操作Service
* @createDate 2024-09-19 16:16:28
*/
public interface OperateLogService extends IService<OperateLog> {

    PageResult getTimelinePage(TimelinePageDto timelinePageDto);
}
