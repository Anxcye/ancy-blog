package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ProjectCardVo;
import com.anxcye.domain.vo.ProjectDetailVo;
import com.anxcye.service.ProjectService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/project")
public class ProjectController {
    @Autowired
    private ProjectService projectService;

    @GetMapping("/list")
    public ResponseResult<List<ProjectCardVo>> projectList() {
        return ResponseResult.success(projectService.getProjectList());
    }

    @GetMapping("/{id}")
    public ResponseResult<ProjectDetailVo> projectDetail(@PathVariable Long id) {
        return ResponseResult.success(projectService.getProjectDetail(id));
    }
}
