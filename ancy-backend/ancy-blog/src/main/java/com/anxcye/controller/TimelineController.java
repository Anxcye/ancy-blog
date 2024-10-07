package com.anxcye.controller;

import com.anxcye.domain.dto.TimelinePageDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.OperateLogService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/timeline")
public class TimelineController {

    @Autowired
    private OperateLogService operateLogService;

   @GetMapping("/page")
   public ResponseResult<PageResult> timelinePage(@ParameterObject TimelinePageDto timelinePageDto) {
       return ResponseResult.success(operateLogService.getTimelinePage(timelinePageDto));
   }
    
}
