package com.anxcye.runner;

import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

@Component
public class ViewCountRunner implements CommandLineRunner {

    @Autowired
    private ArticleService articleService;


    @Override
    public void run(String... args) throws Exception {
        articleService.initViewCount();
    }
}
