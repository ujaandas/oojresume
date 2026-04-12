{
  enable = true;
  name = "research";

  identity = {
    name = "Ujaan Das";
    email = "ujaandas03@gmail.com";
    linkedin = "linkedin.com/in/ujaandas";
    github = "github.com/ujaandas";
  };

  sections = [
    {
      title = "Education";
      entries = [
        "edu_warwick"
        "edu_hkust"
        "edu_northwestern"
      ];
      entryVSpace = 0;
      sectionVSpace = -8;
    }
    {
      title = "Experience";
      entries = [
        "work_hkust_castle_2024"
        "work_stellerus_swe_2025"
      ];
      sectionVSpace = -8;
    }
    {
      title = "Projects";
      entries = [
        "proj_dissertation"
        "proj_follow_me_robot"
      ];
      sectionVSpace = -8;
    }
    {
      title = "Skills";
      entries = [ "skills_default" ];
      sectionVSpace = -8;
    }
  ];
}
