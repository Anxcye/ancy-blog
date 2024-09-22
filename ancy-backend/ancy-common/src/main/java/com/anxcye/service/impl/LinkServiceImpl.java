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
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author axy
 * @description 针对表【ancy_link(友链)】的数据库操作Service实现
 * @createDate 2024-09-10 20:33:41
 */
@Service
public class LinkServiceImpl extends ServiceImpl<LinkMapper, Link>
        implements LinkService {

    @Log
    @Override
    public List<LinkVo> getApprovedLinks() {
        LambdaQueryWrapper<Link> linkLambdaQueryWrapper = new LambdaQueryWrapper<>();
        linkLambdaQueryWrapper.eq(Link::getStatus, SystemConstants.LINK_STATUS_NORMAL);
        List<Link> list = this.list(linkLambdaQueryWrapper);

        List<LinkVo> linkVos = BeanCopyUtils.copyList(list, LinkVo.class);

        return linkVos;
    }

    @Override
    public PageResult pageList(LinkListDto linkListDto) {
        LambdaQueryWrapper<Link> linkLambdaQueryWrapper = new LambdaQueryWrapper<>();

        Page<Link> page = new Page<>(linkListDto.getPageNum(), linkListDto.getPageSize());
        page(page, linkLambdaQueryWrapper);
        return new PageResult(page.getTotal(), page.getRecords());
    }

    @Override
    public boolean addLink(LinkDto linkDto) {
        Link link = BeanCopyUtils.copyBean(linkDto, Link.class);
        save(link);
        return true;
    }

    @Override
    public boolean deleteLink(Long id) {
        removeById(id);
        return true;
    }

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
}




