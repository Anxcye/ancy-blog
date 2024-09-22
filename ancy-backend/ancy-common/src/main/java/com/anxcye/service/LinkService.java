package com.anxcye.service;

import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.dto.LinkListDto;
import com.anxcye.domain.entity.Link;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.LinkVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【ancy_link(友链)】的数据库操作Service
* @createDate 2024-09-10 20:33:41
*/
public interface LinkService extends IService<Link> {

    List<LinkVo> getApprovedLinks();

    PageResult pageList(LinkListDto linkListDto);

    boolean addLink(LinkDto linkDto);

    boolean deleteLink(Long id);

    boolean updateLink(Long id, LinkDto linkDto);

    LinkVo getLink(Long id);
}
