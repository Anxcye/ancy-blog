package com.anxcye.controller;

import com.anxcye.domain.dto.ProjectDto;
import com.anxcye.domain.dto.ProjectPageDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ProjectService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;
@RestController
@RequestMapping("/project")
public class ProjectController {

    @Autowired
    private ProjectService projectService;

    @PreAuthorize("@ps.hasPermission('content:project:list')")
    @GetMapping
    public ResponseResult projectPage(@ParameterObject ProjectPageDto projectPageDto) {
        return ResponseResult.success(projectService.getProjectPage(projectPageDto));
    }
    
    @PreAuthorize("@ps.hasPermission('content:project:add')")
    @PostMapping
    public ResponseResult<Long> projectAdd(@RequestBody ProjectDto projectDto) {
        return ResponseResult.success(projectService.addProject(projectDto));
    }

    @PreAuthorize("@ps.hasPermission('content:project:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> projectUpdate(@PathVariable Long id, @RequestBody ProjectDto projectDto) {
        return ResponseResult.success(projectService.updateProject(id, projectDto));
    }

    @PreAuthorize("@ps.hasPermission('content:project:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> projectDelete(@PathVariable Long id) {
        return ResponseResult.success(projectService.deleteProject(id));
    }
}
