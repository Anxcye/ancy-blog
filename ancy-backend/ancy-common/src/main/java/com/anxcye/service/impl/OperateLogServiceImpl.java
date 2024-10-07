package com.anxcye.service.impl;

import com.anxcye.domain.dto.TimelinePageDto;
import com.anxcye.domain.entity.OperateLog;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.TimelineVo;
import com.anxcye.mapper.OperateLogMapper;
import com.anxcye.service.OperateLogService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.List;

/**
* @author axy
* @description 针对表【sys_operate_log(操作日志表)】的数据库操作Service实现
* @createDate 2024-09-19 16:16:28
*/
@Service
public class OperateLogServiceImpl extends ServiceImpl<OperateLogMapper, OperateLog>
    implements OperateLogService{

    private static final List<String> OPERATE_TYPE_LIST = Arrays.asList("addArticle", "addNote", "addProject");

    @Override
    public PageResult getTimelinePage(TimelinePageDto timelinePageDto) {
        LambdaQueryWrapper<OperateLog> wrapper = new LambdaQueryWrapper<>();
        wrapper.orderByDesc(OperateLog::getOperateTime);
        wrapper.in(OperateLog::getMethodName, OPERATE_TYPE_LIST);

        Page<OperateLog> page = new Page<>(timelinePageDto.getPageNum(), timelinePageDto.getPageSize());
        page(page, wrapper);

        List<TimelineVo> timelineVos = BeanCopyUtils.copyList(page.getRecords(), TimelineVo.class);

        return new PageResult(page.getTotal(), timelineVos);
    }
}




