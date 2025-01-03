package com.anxcye.service;

import com.anxcye.domain.dto.ProjectDto;
import com.anxcye.domain.dto.ProjectPageDto;
import com.anxcye.domain.entity.Project;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ProjectCardVo;
import com.anxcye.domain.vo.ProjectDetailVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【ancy_project(project表)】的数据库操作Service
* @createDate 2024-10-06 20:51:26
*/
public interface ProjectService extends IService<Project> {

    List<ProjectCardVo> getProjectList();

    ProjectDetailVo getProjectDetail(Long id);

    PageResult getProjectPage(ProjectPageDto projectPageDto);

    Long addProject(ProjectDto projectDto);

    Boolean updateProject(Long id, ProjectDto projectDto);

    Boolean deleteProject(Long id);
}
