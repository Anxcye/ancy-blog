package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.ProjectDto;
import com.anxcye.domain.dto.ProjectPageDto;
import com.anxcye.domain.entity.Project;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ProjectCardVo;
import com.anxcye.domain.vo.ProjectDetailVo;
import com.anxcye.mapper.ProjectMapper;
import com.anxcye.service.ArticleService;
import com.anxcye.service.ProjectService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author axy
 * @description 针对表【ancy_project(project表)】的数据库操作Service实现
 * @createDate 2024-10-06 20:51:26
 */
@Service
public class ProjectServiceImpl extends ServiceImpl<ProjectMapper, Project>
        implements ProjectService {

    @Autowired
    private ArticleService articleService;

    private LambdaQueryWrapper<Project> getProjectWrapper() {
        LambdaQueryWrapper<Project> wrapper = new LambdaQueryWrapper<>();
        if (!SecurityUtil.isAdmin()) {
            wrapper.eq(Project::getStatus, SystemConstants.STATUS_NORMAL);
        }
        return wrapper;
    }

    @Override
    public List<ProjectCardVo> getProjectList() {
        LambdaQueryWrapper<Project> wrapper = getProjectWrapper();
        wrapper.orderByDesc(Project::getIsTop)
                .orderByAsc(Project::getOrderNum)
                .orderByDesc(Project::getCreateTime);
        List<Project> projects = list(wrapper);
        return BeanCopyUtils.copyList(projects, ProjectCardVo.class);
    }

    @Override
    public ProjectDetailVo getProjectDetail(Long id) {
        LambdaQueryWrapper<Project> wrapper = getProjectWrapper();
        wrapper.eq(Project::getId, id);
        Project project = getOne(wrapper);

        return BeanCopyUtils.copyBean(project, ProjectDetailVo.class);
    }

    @Override
    public PageResult getProjectPage(ProjectPageDto projectPageDto) {
        LambdaQueryWrapper<Project> wrapper = getProjectWrapper();
        wrapper.orderByDesc(Project::getIsTop)
                .orderByAsc(Project::getOrderNum)
                .orderByDesc(Project::getCreateTime);
        Page<Project> page = new Page<>(projectPageDto.getPageNum(), projectPageDto.getPageSize());
        page(page, wrapper);

        List<ProjectCardVo> projectCardVos = BeanCopyUtils.copyList(page.getRecords(), ProjectCardVo.class);

        return new PageResult(page.getTotal(), projectCardVos);
    }

    @Override
    public Long addProject(ProjectDto projectDto) {
        Project project = BeanCopyUtils.copyBean(projectDto, Project.class);
        save(project);
        return project.getId();
    }

    @Override
    public Boolean updateProject(Long id, ProjectDto projectDto) {
        Project project = BeanCopyUtils.copyBean(projectDto, Project.class);
        project.setId(id);
        return updateById(project);
    }

    @Override
    public Boolean deleteProject(Long id) {
        return removeById(id);
    }
}
