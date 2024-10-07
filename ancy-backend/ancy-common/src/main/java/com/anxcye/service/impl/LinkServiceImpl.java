package com.anxcye.service.impl;

import com.anxcye.annotation.Log;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.dto.LinkListDto;
import com.anxcye.domain.entity.Link;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.LinkVo;
import com.anxcye.mapper.LinkMapper;
import com.anxcye.service.LinkService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.toolkit.StringUtils;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;

/**
 * @author axy
 * @description 针对表【ancy_link(友链)】的数据库操作Service实现
 * @createDate 2024-09-10 20:33:41
 */
@Service
public class LinkServiceImpl extends ServiceImpl<LinkMapper, Link>
        implements LinkService {

    private LambdaQueryWrapper<Link> getLinkWrapper() {
        LambdaQueryWrapper<Link> wrapper = new LambdaQueryWrapper<>();
        if (!SecurityUtil.isAdmin()) {
            wrapper.eq(Link::getStatus, SystemConstants.LINK_STATUS_NORMAL);
        }
        return wrapper;
    }

    @Override
    public List<LinkVo> getApprovedLinks() {
        LambdaQueryWrapper<Link> linkLambdaQueryWrapper = getLinkWrapper();
        List<Link> list = this.list(linkLambdaQueryWrapper);

        return BeanCopyUtils.copyList(list, LinkVo.class);
    }

    @Override
    public PageResult pageList(LinkListDto linkListDto) {
        LambdaQueryWrapper<Link> linkLambdaQueryWrapper = getLinkWrapper();

        linkLambdaQueryWrapper.like(StringUtils.isNotEmpty(linkListDto.getName()),
                Link::getName, linkListDto.getName());
        linkLambdaQueryWrapper.eq(Objects.nonNull(linkListDto.getStatus()),
                Link::getStatus, linkListDto.getStatus());

        Page<Link> page = new Page<>(linkListDto.getPageNum(), linkListDto.getPageSize());
        page(page, linkLambdaQueryWrapper);
        return new PageResult(page.getTotal(), page.getRecords());
    }

    @Log
    @Override
    public Long addLink(LinkDto linkDto) {
        Link link = BeanCopyUtils.copyBean(linkDto, Link.class);
        save(link);
        return link.getId();
    }

    @Log
    @Override
    public boolean deleteLink(Long id) {
        removeById(id);
        return true;
    }

    @Log
    @Override
    public boolean updateLink(Long id, LinkDto linkDto) {
        Link link = BeanCopyUtils.copyBean(linkDto, Link.class);
        link.setId(id);
        updateById(link);
        return true;
    }

    @Override
    public LinkVo getLink(Long id) {
        Link link = getById(id);
        return BeanCopyUtils.copyBean(link, LinkVo.class);
    }

    @Log
    @Override
    public Long addLinkBlog(LinkDto linkDto) {
        Link link = BeanCopyUtils.copyBean(linkDto, Link.class);
        link.setStatus(SystemConstants.LINK_STATUS_HIDE);
        save(link);
        return link.getId();
    }
}
