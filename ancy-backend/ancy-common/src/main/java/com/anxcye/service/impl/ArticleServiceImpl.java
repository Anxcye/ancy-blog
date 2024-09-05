package com.anxcye.service.impl;


import com.anxcye.domin.entity.Article;
import com.anxcye.mapper.ArticleMapper;
import com.anxcye.service.ArticleService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

/**
 * @author axy
 * @description 针对表【ancy_article(文章表)】的数据库操作Service实现
 * @createDate 2024-09-05 11:29:50
 */
@Service
public class ArticleServiceImpl extends ServiceImpl<ArticleMapper, Article>
        implements ArticleService {

}
