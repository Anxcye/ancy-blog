package com.anxcye.service.impl;

import com.anxcye.annotation.Log;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.vo.LinkVo;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Link;
import com.anxcye.service.LinkService;
import com.anxcye.mapper.LinkMapper;
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
}




